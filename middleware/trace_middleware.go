package middleware

import (
	"context"

	"github.com/everfir/logger-go"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
)

func TraceMiddleware(c *gin.Context) {
	if c.Request.Context().Value("span") != nil {
		c.Next()
		return
	}

	ctx := logger.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
	ctx, span := logger.Start(ctx, "tracingMiddleware")
	defer span.End()

	ctx = context.WithValue(ctx, "span", span)
	ctx = context.WithValue(ctx, "baggage", baggage.FromContext(ctx))
	c.Request = c.Request.WithContext(ctx)
	c.Set("span", span)
	c.Set("baggage", baggage.FromContext(ctx))
	c.Next()
}
