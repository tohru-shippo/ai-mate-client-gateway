package app

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	authclient "github.com/tohru-shippo/ai-mate-client-gateway/internal/client/auth"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/handler"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/middleware"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// 装配用户端 HTTP 服务、基础路由和 Server gRPC 客户端。
func NewServer(cfg *Config, logger *slog.Logger) (*http.Server, error) {
	if cfg.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	coreClient, err := authclient.New(cfg.Server.GRPCTarget, cfg.Server.Timeout)
	if err != nil {
		return nil, err
	}

	router := gin.New()
	router.Use(middleware.Recovery(logger))
	router.Use(otelgin.Middleware(cfg.AppName))
	router.Use(middleware.TraceIDHeader())
	router.Use(middleware.Logging(logger))

	api := router.Group("/api")
	handler.NewHealthHandler(cfg.AppName, cfg.Env).Register(api)
	authService := &service.AuthService{Server: coreClient}
	authHandler := &handler.AuthHandler{Service: authService}
	authHandler.Register(api)
	protected := api.Group("")
	protected.Use(middleware.Auth(authService))
	protected.POST("/auth/logout", authHandler.Logout)
	protected.GET("/auth/me", authHandler.Me)

	// Swagger 文档 UI，访问 /swagger/index.html。
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, 404, apperrors.MessageForCode(404, c.GetHeader("X-Language")))
	})
	router.NoMethod(func(c *gin.Context) {
		response.Error(c, http.StatusMethodNotAllowed, 405, apperrors.MessageForCode(405, c.GetHeader("X-Language")))
	})

	server := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}
	server.RegisterOnShutdown(func() {
		if err := coreClient.Close(); err != nil {
			logger.Error("close server grpc failed", slog.Any("error", err))
		}
	})
	return server, nil
}
