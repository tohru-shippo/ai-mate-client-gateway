package service

import (
	"context"

	coreclient "github.com/tohru-shippo/ai-mate-client-gateway/internal/client/core"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/dto"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/middleware"
)

// 用户登录、登出和 token 校验入口逻辑。
type AuthService struct {
	// ai-mate-core gRPC 客户端。
	Core *coreclient.Client
}

// 校验用户登录请求并返回访问令牌。
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	return nil, apperrors.NotImplemented()
}

// 注销当前用户访问令牌。
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return apperrors.NotImplemented()
}

// 校验访问令牌并返回用户身份。
func (s *AuthService) Authenticate(ctx context.Context, token string) (*middleware.UserPrincipal, error) {
	return nil, apperrors.NotImplemented()
}
