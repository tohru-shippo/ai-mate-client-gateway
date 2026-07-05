package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
)

// 用户端网关健康检查接口。
type HealthHandler struct {
	// 服务名。
	appName string
	// 运行环境。
	env string
	// 服务启动时间。
	started time.Time
}

// 创建带启动时间的健康检查处理器。
func NewHealthHandler(appName string, env string) *HealthHandler {
	return &HealthHandler{appName: appName, env: env, started: time.Now()}
}

// 注册 GET /api/health，接口无需鉴权，返回服务基础状态。
func (h *HealthHandler) Register(r gin.IRouter) {
	r.GET("/health", h.Health)
}

// 健康检查响应数据。
type HealthResponse struct {
	// 服务名。
	AppName string `json:"appName"`
	// 运行环境。
	Env string `json:"env"`
	// 服务状态。
	Status string `json:"status"`
	// 服务启动时间。
	StartedAt string `json:"startedAt"`
}

// @Summary      健康检查
// @Description  返回服务名、环境、状态与启动时间
// @Tags         health
// @Produce      json
// @Success      200 {object} response.Body{data=HealthResponse}
// @Router       /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, HealthResponse{
		AppName:   h.appName,
		Env:       h.env,
		Status:    "ok",
		StartedAt: h.started.Format(time.RFC3339),
	})
}
