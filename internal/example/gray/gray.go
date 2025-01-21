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

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/gray"
	"github.com/everfir/go-helpers/middleware"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/gin-gonic/gin"
)

func main() {
	server := InitServer()
	MockClient(server)
}

func InitServer() *http.Server {
	// 创建Gin引擎
	router := gin.Default()

	// 添加灰度中间件
	router.Use(middleware.GrayMiddleware)

	// 设置路由
	router.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()
		accountID := c.GetHeader("account_id")

		if gray.Gray(ctx, "feature_test", accountID) {
			logger.Info(
				ctx,
				"灰度功能已开启",
				field.String("feature", "feature_test"),
				field.String("user", accountID),
				field.String("business", env.Business(ctx)),
			)

			c.String(http.StatusOK, "灰度功能已开启，允许访问")
			return
		}

		logger.Info(
			ctx,
			"灰度功能未开启",
			field.String("feature", "feature_test"),
			field.String("user", accountID),
			field.String("business", env.Business(ctx)),
		)
		c.String(http.StatusUnauthorized, "灰度功能未开启，禁止访问")
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

func MockClient(server *http.Server) {
	// 模拟客户端请求
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			makeRequest(i)
			time.Sleep(1 * time.Second)
		}
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

func makeRequest(i int) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:10083", nil)
	if err != nil {
		fmt.Printf("\x1b[31m创建请求 %d 失败: %v\x1b[0m\n", i, err)
		return
	}

	req.Header.Add(env.BusinessKey.String(), "momo")
	req.Header.Add("account_id", fmt.Sprintf("%d", i))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\x1b[31m请求 %d 失败: %v\x1b[0m\n", i, err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("\x1b[32m请求 %d 成功: %s\x1b[0m\n", i, resp.Status)
	case http.StatusUnauthorized:
		fmt.Printf("\x1b[33m请求 %d: 服务不可用 (灰度未开启)\x1b[0m\n", i)
	default:
		fmt.Printf("\x1b[31m请求 %d: 未知状态码 %d\x1b[0m\n", i, resp.StatusCode)
	}
}
