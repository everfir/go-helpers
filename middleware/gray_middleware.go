package middleware

import (
	"context"

	"github.com/everfir/go-helpers/env"
	"github.com/gin-gonic/gin"
)

func GrayMiddleware(c *gin.Context) {
	// 根据header中的字段来确定业务
	business := c.GetHeader(env.BusinessKey.String())
	if business == "" {
		c.Next()
		return
	}

	// business info
	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, env.BusinessKey, business)
	c.Set(env.BusinessKey.String(), business)
	c.Request = c.Request.WithContext(ctx)

	// ab分组信息, 由apiGateway填充
	router := c.GetHeader(env.RouterKey.String())
	if router != "" {
		c.Writer.Header().Add(env.RouterGroupKey.String(), router)
	}

	c.Next()
}
