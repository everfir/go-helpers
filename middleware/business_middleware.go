package middleware

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs"
)

var getBusinessConfig func() *define.Config[structs.BusinessConfig] = sync.OnceValue(func() *define.Config[structs.BusinessConfig] {
	config, err := nacos.GetConfigFromNacosAndConfigOnChange[structs.BusinessConfig](nacos.GetNacosClient(), "business.json")
	if err != nil {
		panic(err.Error())
	}
	return config
})

func BusinessMiddleware(c *gin.Context) {
	// 根据header中的字段来确定业务
	business := c.GetHeader(env.BusinessKey.String())
	valid := getBusinessConfig().Get().Valid(business)
	if !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err_code": http.StatusBadRequest,
			"err_msg":  "business field in header is not expected, this business does not exist",
		})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, env.BusinessKey, business)
	c.Set(env.BusinessKey.String(), business)
	c.Request = c.Request.WithContext(ctx)

	c.Next()

}
