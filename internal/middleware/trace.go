package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

// 把当前 trace ID 写入 X-Trace-Id 响应头，供客户端关联服务端日志。
// 必须注册在 otelgin 中间件之后，确保 span 已创建。
func TraceIDHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		spanContext := trace.SpanContextFromContext(c.Request.Context())
		if spanContext.HasTraceID() {
			c.Writer.Header().Set("X-Trace-Id", spanContext.TraceID().String())
		}
		c.Next()
	}
}
