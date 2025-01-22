package middleware

import (
	"sync"

	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/gin-gonic/gin"
)

var shutdownConfig func() *define.Config[map[string]bool] = sync.OnceValue(func() *define.Config[map[string]bool] {
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[map[string]bool](nacos.GetNacosClient(), "shutdown.json")
	if err != nil {
		panic(err.Error())
	}
	return config
})

func ShutdownMiddleware(c *gin.Context) {
	// 根据header中的字段来确定业务
	business := c.GetHeader(env.BusinessKey.String())
	if business == "" {
		c.Next()
		return
	}

	if shutdownConfig().Get()[business] {
		c.AbortWithStatus(599)
		return
	}
	c.Next()
}
