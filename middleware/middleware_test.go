package middleware_test

import (
	"testing"

	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/middleware"
	"github.com/gin-gonic/gin"
)

func TestBusinessMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.BusinessMiddleware)

	router.GET("/test", func(c *gin.Context) {
		t.Logf("business: %s", env.Business(c.Request.Context()))
		t.Logf("platform: %s", env.Platform(c.Request.Context()))
		t.Logf("version: %s", env.Version(c.Request.Context()))
		t.Logf("device: %s", env.Device(c.Request.Context()))
		t.Logf("appType: %s", env.AppType(c.Request.Context()))
	})

	router.Run(":8080")
}
