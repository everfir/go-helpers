package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/gray"
	"github.com/everfir/go-helpers/internal/external_api"
	"github.com/everfir/go-helpers/middleware"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
)

func main() {
	server := InitServer()
	MockClient(server)
}

/*
	测试逻辑：分两组用户测试，先测试A组，在测试B组, 从feature_test1 -> feature_test_10 测试配置如下
	其中：
		- 1、2、3、9为实验中的feature
		- 1、2、4、9为不可靠分支，通过开关控制是否生效, 其中1、4已经回滚， 2、9还在灰度中
	预期结果：
		A组: 1, 3, 4, 9 不可用, 其他可用
		B组: 1, 2, 4 不可用, 其他可用


	// 测试配置如下，从feature_test_1 -> feature_test_9, 预期结果：
	//
	{
	 	// 针对helper_test开启灰度配置
		"helper_test": {
			"experiment": {
				"rule": {
					"mode": 2,	// 混合分流模式(0: 比例模式, 2: 白名单模式, 3: 混合模式)
					"rate": 0.5,	// 1:1分流
					"whitelist": [	// 白名单列表
						"user1",
						"user2"
					]
				},

				// 可进行AB实验的feature列表
				"feature_list": [
					"feature_test_1",
					"feature_test_2",
					"feature_test_3",
					"feature_test_9"
				],

				// A组实验feature列表
				"experiments_a": [
					"feature_test_1",
					"feature_test_2"
				],

				// B组实验feature列表
				"experiments_b": [
					"feature_test_3",
					"feature_test_9"
				]
			},

			// feature灰度功能开关
			"feature_gray": {
				"feature_test_1": false,
				"feature_test_2": true,
				"feature_test_4": false,
				"feature_test_9": true
			}
		}
	}

*/

func InitServer() *http.Server {
	// 创建Gin引擎
	router := gin.New()

	// 添加灰度中间件
	router.Use(middleware.BaseMiddlewares()...)

	// 设置路由
	router.GET("/feature/:name", func(c *gin.Context) {
		ctx := c.Request.Context()
		feature := c.Param("name")
		group := env.ExperimentGroup(ctx)

		// 检查灰度状态
		if gray.Gray(ctx, feature) {
			logger.Info(
				ctx,
				"灰度功能已开启",
				field.String("feature", feature),
				field.String("group", group),
				field.String("business", env.Business(ctx)),
			)

			c.String(http.StatusOK, fmt.Sprintf("功能 %s 已开启，允许访问", feature))
			return
		}

		logger.Info(
			ctx,
			"灰度功能未开启",
			field.String("feature", feature),
			field.String("group", group),
			field.String("business", env.Business(ctx)),
		)
		c.String(http.StatusUnauthorized, fmt.Sprintf("功能 %s 未开启，禁止访问", feature))
	})

	// 创建HTTP服务
	server := &http.Server{
		Addr:    ":10083",
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Println("服务器正在启动，监听端口 :10083")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(context.TODO(), "服务器启动失败", field.String("err", err.Error()))
		}
	}()

	return server
}

// MockClient 模拟客户端请求
func MockClient(server *http.Server) {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// 测试A组用户
		makeRequest("", "A")

		// 测试B组用户
		makeRequest("", "B")

		// 测试routerKey
		makeRequest("test key", "")

	}()

	// 处理优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal(context.TODO(), "服务器关闭失败", field.String("err", err.Error()))
	}

	wg.Wait()
	log.Println("服务器已成功关闭")
}

// makeRequest 发送HTTP请求
func makeRequest(routerKey string, group string) {
	features := []string{
		"feature_test_1",
		"feature_test_2",
		"feature_test_3",
		"feature_test_4",
		"feature_test_5",
		"feature_test_6",
		"feature_test_7",
		"feature_test_8",
		"feature_test_9",
		"feature_test_10",
	}

	succFeatures := []string{}
	failFeatures := []string{}

	userID := ""
	for _, feature := range features {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:10083/feature/"+feature, nil)
		if err != nil {
			fmt.Printf("\x1b[31m[%s] 创建请求失败: %v\x1b[0m\n", userID, err)
			continue
		}

		req.Header.Add(env.BusinessKey.String(), "helper_test")
		token := external_api.GetAccountCfg().Get().GuestToken
		req.Header.Add(env.Authorization.String(), token)
		// RouterKey不为空则测试根据RouterKey获取实验组，group不为空则指定实验组进行测试
		if routerKey != "" {
			req.Header.Add(env.RouterKey.String(), routerKey)
		}
		if group != "" {
			req.Header.Add(env.RouterGroupKey.String(), group)
		}
		logger.Info(req.Context(), "发送请求", field.String("feature", feature), field.String("group", group),
			field.String("routerKey", routerKey), field.String("token", token))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("\x1b[31m[%s] 请求 %s 失败: %v\x1b[0m\n", userID, feature, err)
			continue
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			succFeatures = append(succFeatures, feature)
			// fmt.Printf("\x1b[32m[%s] 请求 %s 成功: %s\x1b[0m\n", userID, feature, resp.Status)
		case http.StatusUnauthorized:
			failFeatures = append(failFeatures, feature)
			// fmt.Printf("\x1b[33m[%s] 请求 %s: 服务不可用 (灰度未开启)\x1b[0m\n", userID, feature)
		default:
			fmt.Printf("\x1b[31m[%s] 请求 %s: 未知状态码 %d\x1b[0m\n", userID, feature, resp.StatusCode)
		}
	}

	logger.Info(
		context.TODO(),
		"灰度测试结果",
		field.String("group", group),
		field.Any("succ", succFeatures),
		field.Any("fail", failFeatures),
	)
}
