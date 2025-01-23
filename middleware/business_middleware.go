package middleware

import (
	"context"
	"net/http"

	"github.com/everfir/go-helpers/env"
	"github.com/gin-gonic/gin"
)

func BusinessMiddleware(c *gin.Context) {
	// 根据header中的字段来确定业务
	business := c.GetHeader(env.BusinessKey.String())
	if business == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"err_code": http.StatusBadRequest,
			"err_msg":  "unexcept business field in header",
		})
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, env.BusinessKey, business)
	c.Set(env.BusinessKey.String(), business)
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}
