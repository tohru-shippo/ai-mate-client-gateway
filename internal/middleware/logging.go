package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// 记录用户端 HTTP 请求访问日志。
func Logging(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		spanContext := trace.SpanContextFromContext(c.Request.Context())
		attrs := []slog.Attr{
			slog.String("requestId", GetRequestID(c)),
			slog.String("method", c.Request.Method),
			slog.String("path", c.FullPath()),
			slog.Int("status", c.Writer.Status()),
			slog.Int64("latencyMs", time.Since(start).Milliseconds()),
			slog.String("clientIp", c.ClientIP()),
		}
		if spanContext.HasTraceID() {
			attrs = append(attrs, slog.String("traceId", spanContext.TraceID().String()), slog.String("spanId", spanContext.SpanID().String()))
		}

		logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "http request completed", attrs...)
	}
}
