package middleware

import (
	"context"

	"github.com/everfir/go-helpers/env"
	"github.com/gin-gonic/gin"
)

func GrayMiddleware(c *gin.Context) {
	group := c.GetHeader(env.RouterGroupKey.String())

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, env.ExperimentGroupKey, group)
	c.Set(env.RouterGroupKey.String(), group)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
