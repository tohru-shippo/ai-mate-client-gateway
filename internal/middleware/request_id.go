package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Gin 上下文中保存请求 ID 的键。
const RequestIDKey = "requestId"

// 为每个请求生成或透传 X-Request-Id。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.Must(uuid.NewV7()).String()
		}
		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set("X-Request-Id", requestID)
		c.Next()
	}
}

// 从 Gin 上下文读取请求 ID。
func GetRequestID(c *gin.Context) string {
	value, ok := c.Get(RequestIDKey)
	if !ok {
		return ""
	}
	requestID, _ := value.(string)
	return requestID
}
