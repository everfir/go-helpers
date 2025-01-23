package middleware

import (
	"net/http"
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
	business := env.Business(c.Request.Context())
	if business == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err_code": http.StatusBadRequest,
			"err_msg":  "unexcept business field in header",
		})
		return
	}

	if shutdownConfig().Get()[business] {
		c.AbortWithStatus(599)
		return
	}
	c.Next()
}
