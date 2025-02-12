package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/gray"
	"github.com/everfir/go-helpers/middleware"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/logger-go/structs/log_level"
	"github.com/gin-gonic/gin"
)

func main() {
	// 注释开启info等级日志
	logger.Init(logger.WithLevel(log_level.ErrorLevel))

	server := InitServer()
	MockClient(server)
}

/*
		测试逻辑： 待补充
		测试用户：
		 第一组: 满足条件为[1,2,3,4]
			users := [][]string{
				{"1", "android"},
				{"2", "android"},
				{"3", "ios"},
				{"4", "ios"},
				{"5", "android"},
			}

		 第二组: 满足条件为[12, 13,564]
			users := [][]string{
				{"1", "android"},
				{"12", "ios"},
				{"13", "ios"},
				{"123", "android"},
				{"564", "ios"},
			}
		// 测试配置如下，预期结果：
		- 第一组条件，满足条件为
		//
		{
	    "helper_test": {
	        "feature_test": {
	            "enable": true,
	            "rule": [
	                {
	                    "mode": 4,	// 平台分流
	                    "targets": [
	                        "ios"
	                    ]
	                },
	                {
	                    "mode": 1,	// 规则分流
	                    "expresion": "user.account_id % 2 == 0"
	                }
	            ]
	        }
	    }
	}
*/
func InitServer() *http.Server {
	// 创建Gin引擎
	router := gin.New()

	// 添加灰度中间件
	router.Use(middleware.BaseMiddlewares()...)
	router.Use(middleware.AuthMiddleware())

	// 设置路由
	router.GET("/feature", func(c *gin.Context) {
		ctx := c.Request.Context()
		feature := "feature_test"

		id, err := strconv.Atoi(c.GetHeader("account_id"))
		if err != nil {
			logger.Error(
				ctx,
				"invalid account_id",
				field.String("id", c.GetHeader("account_id")),
			)
		}
		ctx = env.TestSetAccountInfo(ctx, uint64(id))

		logger.Error(
			ctx,
			"info",
			field.String("device", env.Device(ctx).String()),
			field.String("platform", env.Platform(ctx).String()),
			field.String("version", env.Version(ctx)),
			field.Any("account", env.AccountInfo(ctx)),
		)

		// 检查灰度状态
		if gray.Experimental(ctx, feature) {
			logger.Info(
				ctx,
				"灰度功能已开启",
				field.String("feature", feature),
				field.String("business", env.Business(ctx)),
			)

			c.String(http.StatusOK, fmt.Sprintf("功能 %s 已开启，允许访问", feature))
			return
		}

		logger.Info(
			ctx,
			"灰度功能未开启",
			field.String("feature", feature),
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

		users := [][]string{
			{"1", "android"},
			{"2", "android"},
			{"3", "ios"},
			{"4", "ios"},
			{"5", "android"},
		}
		// 测试A组用户
		makeRequest(users)

		users = [][]string{
			{"1", "android"},
			{"12", "ios"},
			{"13", "ios"},
			{"123", "android"},
			{"564", "ios"},
		}
		// 测试B组用户
		makeRequest(users)
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
func makeRequest(users [][]string) {

	var succ = []string{}
	var fail = []string{}

	for _, user := range users {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:10083/feature/", nil)
		if err != nil {
			fmt.Printf("\x1b[31m[%s] 创建请求失败: %v\x1b[0m\n", user, err)
			continue
		}

		req.Header.Add("account_id", user[0])
		req.Header.Add(env.BusinessKey.String(), "helper_test")
		req.Header.Add(env.VersionKey.String(), "0.0.10")
		req.Header.Add(env.PlatformKey.String(), string(user[1]))
		req.Header.Add(env.DeviceKey.String(), env.Dev_Phone.String())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("请求 %s 失败: %v\x1b[0m\n", user, err)
			continue
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusOK:
			succ = append(succ, user[0])
			// fmt.Printf("\x1b[32m[%s] 请求 %s 成功: %s\x1b[0m\n", userID, feature, resp.Status)
		case http.StatusUnauthorized:
			fail = append(fail, user[0])
			// fmt.Printf("\x1b[33m[%s] 请求 %s: 服务不可用 (灰度未开启)\x1b[0m\n", userID, feature)
		default:
			fmt.Printf("请求 %s: 未知状态码 %d\x1b[0m\n", user, resp.StatusCode)
		}
	}

	logger.Error(
		context.TODO(),
		"灰度测试结果",
		field.Any("succ", succ),
		field.Any("fail", fail),
	)
}
