package service

import (
	"context"

	authclient "github.com/tohru-shippo/ai-mate-client-gateway/internal/client/auth"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/dto"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/middleware"
)

// 用户登录、登出、刷新和 token 校验入口逻辑。
type AuthService struct {
	// ai-mate-server gRPC 客户端。
	Server *authclient.Client
}

// 校验用户登录请求并返回访问令牌，按 method 分发到密码或 OAuth 登录。
func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	switch req.Method {
	case "account_password":
		return s.Server.SignInWithPassword(ctx, req.Email, req.Password)
	case "google", "twitter":
		return s.Server.SignInWithOAuth(ctx, req)
	default:
		return nil, apperrors.BadRequest("")
	}
}

// 刷新访问令牌。
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	return s.Server.RefreshToken(ctx, refreshToken)
}

// 注销当前用户会话。
func (s *AuthService) Logout(ctx context.Context, userID, sessionID string) error {
	return s.Server.Logout(ctx, userID, sessionID)
}

// 查询当前用户资料。
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*dto.MyProfile, error) {
	return s.Server.GetMyProfile(ctx, userID)
}

// 校验访问令牌并返回用户身份。
func (s *AuthService) Authenticate(ctx context.Context, token string) (*middleware.UserPrincipal, error) {
	userID, userType, status, sessionID, err := s.Server.ValidateAccessToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &middleware.UserPrincipal{
		UserID:    userID,
		UserType:  userType,
		Status:    status,
		SessionID: sessionID,
	}, nil
}
