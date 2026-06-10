package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
)

// 捕获 panic 并返回统一 JSON 错误响应。
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.ErrorContext(c.Request.Context(), "http panic recovered", slog.String("requestId", GetRequestID(c)), slog.Any("panic", recovered))
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.Body{Code: apperrors.InternalError.Code, Message: apperrors.InternalError.Message})
	})
}
