package middleware

import (
	"net/http"

	"github.com/everfir/go-helpers/helper/account"
	"github.com/gin-gonic/gin"
)

const (
	ERR_CODE_INVALID_TOKEN = 2000002
	ERR_CODE_TOKEN_MISSING = 2000004

	ERR_MSG_INVALID_TOKEN = "invalid token"
	ERR_MSG_TOKEN_MISSING = "token missing"
)

// AuthMiddleware 访问 check token 接口，并将用户信息存储在 Context 中
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"err_code": ERR_CODE_TOKEN_MISSING, "err_msg": ERR_MSG_TOKEN_MISSING})
			c.Abort()
			return
		}

		ctx, err := account.CheckToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"err_code": ERR_CODE_INVALID_TOKEN,
				},
			)
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
