package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/everfir/go-helpers/consts"
	"github.com/everfir/go-helpers/define/config"
	"github.com/everfir/go-helpers/internal/helper/nacos"
	. "github.com/everfir/go-helpers/internal/structs"
)

var getBusinessConfig func() *config.NacosConfig[BusinessConfig] = sync.OnceValue(func() *config.NacosConfig[BusinessConfig] {
	config, err := nacos.GetConfigAndListen[BusinessConfig](nacos.GetNacosClient(), "business.json")
	if err != nil {
		panic(err.Error())
	}
	return config
})

func BusinessMiddleware(c *gin.Context) {
	// 根据header中的字段来确定业务信息
	business := strings.ToLower(c.GetHeader(consts.BusinessKey.String()))
	platform := strings.ToLower(c.GetHeader(consts.PlatformKey.String()))
	version := strings.ToLower(c.GetHeader(consts.VersionKey.String()))
	device := strings.ToLower(c.GetHeader(consts.DeviceKey.String()))
	appType := strings.ToLower(c.GetHeader(consts.AppTypeKey.String()))

	valid := getBusinessConfig().Get().Valid(business)
	if !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err_code": http.StatusBadRequest,
			"err_msg":  "business field in header is not expected, this business does not exist",
		})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, consts.BusinessKey, business)
	ctx = context.WithValue(ctx, consts.PlatformKey, platform)
	ctx = context.WithValue(ctx, consts.DeviceKey, device)
	ctx = context.WithValue(ctx, consts.VersionKey, version)
	ctx = context.WithValue(ctx, consts.AppTypeKey, appType)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
