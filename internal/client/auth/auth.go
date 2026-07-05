package auth

import (
	"context"

	"github.com/tohru-shippo/ai-mate-client-gateway/internal/dto"
	serverv1 "github.com/tohru-shippo/ai-mate-server/pkg/gen/server/v1"
)

// C 端默认客户端平台。
const clientPlatform = "web"

// 通过 Core 邮箱密码登录。
func (c *Client) SignInWithPassword(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := serverv1.NewAuthServiceClient(c.conn).SignInWithPassword(ctx, &serverv1.SignInWithPasswordRequest{
		Email:    email,
		Password: password,
		Client:   &serverv1.ClientContext{ClientPlatform: clientPlatform},
	})
	if err != nil {
		return nil, err
	}
	return toLoginResponse(resp.GetAuth()), nil
}

// 通过 Core 第三方 OAuth 登录。
func (c *Client) SignInWithOAuth(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := serverv1.NewAuthServiceClient(c.conn).SignInWithOAuth(ctx, &serverv1.SignInWithOAuthRequest{
		Provider:         req.Method,
		AuthorizationCode: req.OAuthCode,
		RedirectUri:      req.RedirectURI,
		CodeVerifier:     req.CodeVerifier,
		Client:           &serverv1.ClientContext{ClientPlatform: clientPlatform},
		Language:         req.Language,
		Region:           req.Region,
	})
	if err != nil {
		return nil, err
	}
	return toLoginResponse(resp.GetAuth()), nil
}

// 通过 Core 刷新令牌。
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := serverv1.NewAuthServiceClient(c.conn).RefreshToken(ctx, &serverv1.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		return nil, err
	}
	return toLoginResponse(resp.GetAuth()), nil
}

// 通过 Core 校验访问令牌，返回用户身份与会话。
func (c *Client) ValidateAccessToken(ctx context.Context, accessToken string) (userID, userType, status, sessionID string, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := serverv1.NewAuthServiceClient(c.conn).ValidateAccessToken(ctx, &serverv1.ValidateAccessTokenRequest{AccessToken: accessToken})
	if err != nil {
		return "", "", "", "", err
	}
	return resp.GetUserId(), resp.GetUserType(), resp.GetStatus(), resp.GetSessionId(), nil
}

// 通过 Core 注销会话。
func (c *Client) Logout(ctx context.Context, userID, sessionID string) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := serverv1.NewAuthServiceClient(c.conn).Logout(ctx, &serverv1.LogoutRequest{UserId: userID, SessionId: sessionID})
	return err
}

// 通过 Core 查询当前用户资料。
func (c *Client) GetMyProfile(ctx context.Context, userID string) (*dto.MyProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	resp, err := serverv1.NewAuthServiceClient(c.conn).GetMyProfile(ctx, &serverv1.GetMyProfileRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return &dto.MyProfile{
		UserID:      resp.GetUserId(),
		UserType:    resp.GetUserType(),
		Email:       resp.GetEmail(),
		DisplayName: resp.GetDisplayName(),
		AvatarURL:   resp.GetAvatarUrl(),
		Language:    resp.GetLanguage(),
		Region:      resp.GetRegion(),
		Status:      resp.GetStatus(),
	}, nil
}

// toLoginResponse 将 AuthTokens 映射为登录响应。
func toLoginResponse(auth *serverv1.AuthTokens) *dto.LoginResponse {
	if auth == nil {
		return nil
	}
	return &dto.LoginResponse{
		AccessToken:      auth.GetAccessToken(),
		ExpiresIn:        auth.GetAccessTokenExpiresIn(),
		RefreshToken:     auth.GetRefreshToken(),
		RefreshExpiresIn: auth.GetRefreshTokenExpiresIn(),
		UserID:           auth.GetUserId(),
		UserType:         auth.GetUserType(),
		Status:           auth.GetStatus(),
		IsNewUser:        auth.GetIsNewUser(),
	}
}
