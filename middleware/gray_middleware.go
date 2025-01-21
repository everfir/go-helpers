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

	group := c.GetHeader(env.RouterGroupKey.String())

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, env.BusinessKey, business)
	ctx = context.WithValue(ctx, env.ExperimentGroupKey, group)
	c.Set(env.BusinessKey.String(), business)
	c.Set(env.RouterGroupKey.String(), group)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
