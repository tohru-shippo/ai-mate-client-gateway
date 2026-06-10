package handler

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
)

// 将 Core gRPC 错误映射为用户端多语言统一错误响应。
func writeCoreError(c *gin.Context, err error) {
	response.FromError(c, apperrors.FromCoreError(err, c.GetHeader("X-Language")))
}
