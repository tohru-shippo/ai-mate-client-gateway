package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/dto"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/service"
)

// 用户端认证接口。
type AuthHandler struct {
	// 用户认证服务。
	Service *service.AuthService
}

// 注册 POST /api/auth/login 和 POST /api/auth/logout。
func (h *AuthHandler) Register(r gin.IRouter) {
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)
}

// POST /api/auth/login，无需鉴权，后续接入 Core 登录能力。
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FromError(c, apperrors.BadRequest(""))
		return
	}
	result, err := h.Service.Login(c.Request.Context(), req)
	if err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, result)
}

// POST /api/auth/logout，注销当前用户访问令牌。
func (h *AuthHandler) Logout(c *gin.Context) {
	token := strings.TrimSpace(strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer "))
	if err := h.Service.Logout(c.Request.Context(), token); err != nil {
		response.FromError(c, err)
		return
	}
	response.Success(c, nil)
}
