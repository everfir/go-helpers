package middleware

import (
	"sync"

	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/gin-gonic/gin"
)

var shutdownConfig func() *nacos.Config[map[string]bool] = sync.OnceValue(func() *nacos.Config[map[string]bool] {
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[map[string]bool]("shutdown.json")
	if err != nil {
		panic(err.Error())
	}
	return config
})

func ShutdownMiddleware(c *gin.Context) {
	// 根据header中的字段来确定业务
	business := c.GetHeader("x-everfir-business")
	if business == "" {
		c.Next()
		return
	}

	if shutdownConfig().Data[business] {
		c.AbortWithStatus(599)
		return
	}
	c.Next()
}
