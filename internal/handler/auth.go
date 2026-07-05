package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/dto"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/middleware"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/service"
)

// 用户端认证接口。
type AuthHandler struct {
	// 用户认证服务。
	Service *service.AuthService
}

// 注册公开路由（登出与 me 需鉴权，由 protected 组注册）。
func (h *AuthHandler) Register(r gin.IRouter) {
	r.POST("/auth/login", h.Login)
	r.POST("/auth/refresh", h.Refresh)
}

// POST /api/auth/login，用户登录，按 method 分发密码或 OAuth。
//
// @Summary      用户登录
// @Description  使用邮箱密码或 Google/Twitter OAuth 登录，返回访问令牌与刷新令牌
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        X-Language header string false "界面语言（zh-CN/zh-TW/en-US/ja-JP/ko-KR/ru-RU）"
// @Param        body body dto.LoginRequest true "登录参数"
// @Success      200 {object} response.Body{data=dto.LoginResponse}
// @Failure      400 {object} response.Body "请求参数错误"
// @Failure      401 {object} response.Body "未授权"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FromError(c, apperrors.BadRequest(""))
		return
	}
	result, err := h.Service.Login(c.Request.Context(), req)
	if err != nil {
		writeCoreError(c, err)
		return
	}
	response.Success(c, result)
}

// POST /api/auth/refresh，刷新访问令牌。
//
// @Summary      刷新令牌
// @Description  使用刷新令牌换取新的访问令牌与刷新令牌
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body dto.RefreshRequest true "刷新参数"
// @Success      200 {object} response.Body{data=dto.LoginResponse}
// @Failure      400 {object} response.Body "请求参数错误"
// @Failure      401 {object} response.Body "未授权"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FromError(c, apperrors.BadRequest(""))
		return
	}
	result, err := h.Service.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		writeCoreError(c, err)
		return
	}
	response.Success(c, result)
}

// POST /api/auth/logout，注销当前用户会话，需先通过鉴权。
//
// @Summary      用户登出
// @Description  注销当前登录会话
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <token>"
// @Success      200 {object} response.Body
// @Failure      401 {object} response.Body "未授权"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	principal, err := middleware.MustUser(c)
	if err != nil {
		response.FromError(c, apperrors.Unauthorized())
		return
	}
	if err := h.Service.Logout(c.Request.Context(), principal.UserID, principal.SessionID); err != nil {
		writeCoreError(c, err)
		return
	}
	response.Success(c, nil)
}

// GET /api/auth/me，查询当前用户资料，需先通过鉴权。
//
// @Summary      当前用户资料
// @Description  返回当前登录用户的昵称、头像等资料
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization header string true "Bearer <token>"
// @Success      200 {object} response.Body{data=dto.MyProfile}
// @Failure      401 {object} response.Body "未授权"
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	principal, err := middleware.MustUser(c)
	if err != nil {
		response.FromError(c, apperrors.Unauthorized())
		return
	}
	profile, err := h.Service.GetProfile(c.Request.Context(), principal.UserID)
	if err != nil {
		writeCoreError(c, err)
		return
	}
	response.Success(c, profile)
}
