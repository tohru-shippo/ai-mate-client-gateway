package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
)

// 用户端对外统一响应结构。
type Body struct {
	// 业务码，0 表示成功。
	Code int `json:"code"`
	// 响应消息。
	Message string `json:"message"`
	// 响应数据。
	Data any `json:"data,omitempty"`
}

// 写入成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "success", Data: data})
}

// 写入指定 HTTP 状态码和业务错误响应。
func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Body{Code: code, Message: message})
}

// 将业务错误写入统一响应。
func FromError(c *gin.Context, err error) {
	appErr := apperrors.Localize(err, c.GetHeader("X-Language"))
	if appErr == nil {
		return
	}
	c.JSON(appErr.HTTPStatus, Body{Code: appErr.Code, Message: appErr.Message})
}
