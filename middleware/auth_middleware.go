package middleware

import (
	"context"
	"github.com/everfir/go-helpers/env"
	"github.com/everfir/go-helpers/internal/external_api"
	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthMiddleware 访问 check token 接口，并将用户信息存储在 Context 中
func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader(string(env.Authorization))
	routerKey := c.GetHeader(string(env.RouterKey))
	if token == "" {
		logger.Error(c, "token is empty")
		c.JSON(http.StatusUnauthorized, gin.H{"err_code": http.StatusBadRequest, "err_msg": "token is empty"})
		c.Abort()
		return
	}

	// 访问 check token 接口
	accountInfo, routerGroup, err := external_api.CheckToken(c, token, routerKey)
	if err != nil {
		logger.Error(c, "CheckToken failed", field.String("err", err.Error()), field.String("token", token))
		c.JSON(http.StatusUnauthorized, gin.H{"err_code": http.StatusBadRequest, "err_msg": "token is invalid"})
		c.Abort()
		return
	}

	if !accountInfo.Valid {
		logger.Error(c, "accountInfo is invalid", field.String("token", token), field.Any("accountInfo", accountInfo))
		c.JSON(http.StatusOK, gin.H{"err_code": http.StatusBadRequest, "err_msg": "token is expired"})
		c.Abort()
		return
	}

	// 将用户信息存储到 Context 中
	ctx := context.WithValue(c.Request.Context(), "accountInfo", accountInfo.AccountInfo)

	// 将用户的 routerGroup
	ctx = context.WithValue(ctx, "routerGroup", routerGroup)

	logger.Info(ctx, "CheckToken success", field.String("token", token),
		field.String("router key", c.GetHeader(string(env.RouterKey))),
		field.String("router key", c.GetHeader(string(env.RouterGroupKey))),
		field.String("routerGroup", routerGroup))

	c.Request = c.Request.WithContext(ctx)

	// 继续后续的请求处理
	c.Next()
}
