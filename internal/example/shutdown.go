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

	"everfir/go-helpers/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin引擎
	router := gin.Default()

	// 添加shutdown中间件
	router.Use(middleware.ShutdownMiddleware)

	// 设置路由
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is running normally")
	})

	// 创建HTTP服务
	server := &http.Server{
		Addr:    ":10083",
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Println("Starting server on :10083")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

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

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	wg.Wait()
	log.Println("Server stopped")
}

func makeRequest(i int) {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost:10083", nil)
	req.Header.Add("x-everfir-business", "momo")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("\x1b[31mRequest %d failed: %v\x1b[0m\n", i, err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("\x1b[32mRequest %d success: %s\x1b[0m\n", i, resp.Status)
	case 599:
		fmt.Printf("\x1b[33mRequest %d: Service Unavailable\x1b[0m\n", i)
	default:
		fmt.Printf("\x1b[31mRequest %d: Unexpected status %d\x1b[0m\n", i, resp.StatusCode)
	}
}
