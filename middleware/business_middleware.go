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
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
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

	logger.Debug(c.Request.Context(), "business", field.String("business", business))
	logger.Debug(c.Request.Context(), "platform", field.String("platform", platform))
	logger.Debug(c.Request.Context(), "version", field.String("version", version))
	logger.Debug(c.Request.Context(), "device", field.String("device", device))
	logger.Debug(c.Request.Context(), "appType", field.String("appType", appType))

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
	ctx = context.WithValue(ctx, consts.PlatformKey, consts.TDevicePlatform(platform))
	ctx = context.WithValue(ctx, consts.DeviceKey, consts.TDevice(device))
	ctx = context.WithValue(ctx, consts.VersionKey, version)
	ctx = context.WithValue(ctx, consts.AppTypeKey, consts.TAppType(appType))
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
