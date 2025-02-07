package middleware

import "github.com/gin-gonic/gin"

func BaseMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		BusinessMiddleware,
		AuthMiddleware,
		TraceMiddleware,
		GrayMiddleware,
		ShutdownMiddleware,
	}
}
