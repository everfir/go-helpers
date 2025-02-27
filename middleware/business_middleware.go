package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/everfir/go-helpers/define"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	"github.com/everfir/go-helpers/internal/structs"
)

var getBusinessConfig func() *define.Config[structs.BusinessConfig] = sync.OnceValue(func() *define.Config[structs.BusinessConfig] {
	config, err := nacos.GetConfigAndListen[structs.BusinessConfig](nacos.GetNacosClient(), "business.json")
	if err != nil {
		panic(err.Error())
	}
	return config
})

func BusinessMiddleware(c *gin.Context) {
	// 根据header中的字段来确定业务信息
	business := strings.ToLower(c.GetHeader(env.BusinessKey.String()))
	platform := strings.ToLower(c.GetHeader(env.PlatformKey.String()))
	version := strings.ToLower(c.GetHeader(env.VersionKey.String()))
	device := strings.ToLower(c.GetHeader(env.DeviceKey.String()))
	appType := strings.ToLower(c.GetHeader(env.AppTypeKey.String()))

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
	ctx = context.WithValue(ctx, env.PlatformKey, platform)
	ctx = context.WithValue(ctx, env.DeviceKey, device)
	ctx = context.WithValue(ctx, env.VersionKey, version)
	ctx = context.WithValue(ctx, env.AppTypeKey, appType)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
