package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
)

// 捕获 panic 并返回统一 JSON 错误响应。
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		attrs := []slog.Attr{slog.Any("panic", recovered)}
		if spanContext := trace.SpanContextFromContext(c.Request.Context()); spanContext.HasTraceID() {
			attrs = append(attrs, slog.String("traceId", spanContext.TraceID().String()))
		}
		logger.LogAttrs(c.Request.Context(), slog.LevelError, "http panic recovered", attrs...)
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Body{Code: apperrors.InternalError.Code, Message: apperrors.InternalError.Message})
	})
}
