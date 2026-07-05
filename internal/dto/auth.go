package dto

// 用户登录请求，method 决定登录方式。
type LoginRequest struct {
	// 登录方式：account_password、google、twitter。
	Method string `json:"method"`
	// 邮箱，account_password 登录时使用。
	Email string `json:"email"`
	// 密码，account_password 登录时使用。
	Password string `json:"password"`
	// OAuth 授权码，google/twitter 登录时使用。
	OAuthCode string `json:"oauthCode"`
	// OAuth 重定向 URI。
	RedirectURI string `json:"redirectUri"`
	// PKCE code_verifier。
	CodeVerifier string `json:"codeVerifier"`
	// 当前前端语言。
	Language string `json:"language"`
	// 用户地区。
	Region string `json:"region"`
}

// 用户登录响应。
type LoginResponse struct {
	// 访问令牌。
	AccessToken string `json:"accessToken"`
	// 访问令牌有效期，单位秒。
	ExpiresIn int64 `json:"expiresIn"`
	// 刷新令牌。
	RefreshToken string `json:"refreshToken"`
	// 刷新令牌有效期，单位秒。
	RefreshExpiresIn int64 `json:"refreshExpiresIn"`
	// 用户 ID。
	UserID string `json:"userId"`
	// 用户类型：normal、internal。
	UserType string `json:"userType"`
	// 用户状态：normal、restricted、archived。
	Status string `json:"status"`
	// 是否新用户。
	IsNewUser bool `json:"isNewUser"`
}

// 刷新令牌请求。
type RefreshRequest struct {
	// 刷新令牌。
	RefreshToken string `json:"refreshToken"`
}

// 当前用户资料。
type MyProfile struct {
	// 用户 ID。
	UserID string `json:"userId"`
	// 用户类型。
	UserType string `json:"userType"`
	// 登录邮箱。
	Email string `json:"email"`
	// 展示名。
	DisplayName string `json:"displayName"`
	// 头像地址。
	AvatarURL string `json:"avatarUrl"`
	// 首选语言。
	Language string `json:"language"`
	// 地区。
	Region string `json:"region"`
	// 用户状态。
	Status string `json:"status"`
}
