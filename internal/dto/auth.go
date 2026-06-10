package dto

// 用户登录请求。
type LoginRequest struct {
	// 登录方式，例如 password、email_code、google_oauth。
	Method string `json:"method"`
	// 账号，账号密码登录时使用。
	Account string `json:"account"`
	// 手机号，手机验证码登录时使用。
	Phone string `json:"phone"`
	// 邮箱，邮箱验证码登录时使用。
	Email string `json:"email"`
	// 密码，账号密码登录时使用。
	Password string `json:"password"`
	// 验证码，验证码登录时使用。
	Code string `json:"code"`
	// OAuth 授权码。
	OAuthCode string `json:"oauthCode"`
	// OAuth 重定向 URI。
	RedirectURI string `json:"redirectUri"`
	// PKCE code_verifier。
	CodeVerifier string `json:"codeVerifier"`
	// 当前前端语言。
	Language string `json:"language"`
}

// 用户登录响应。
type LoginResponse struct {
	// 访问令牌。
	AccessToken string `json:"accessToken"`
	// 过期时间，单位秒。
	ExpiresIn int64 `json:"expiresIn"`
	// 用户基础信息。
	User User `json:"user"`
}

// 用户端接口暴露的用户基础信息。
type User struct {
	// 用户 ID。
	ID string `json:"id"`
	// 用户名。
	Username string `json:"username"`
	// 昵称。
	Nickname string `json:"nickname"`
	// 头像 URL。
	Avatar string `json:"avatar"`
	// 用户角色。
	Role string `json:"role"`
	// 用户首选语言。
	Language string `json:"language"`
}

// 用户登出请求。
type LogoutRequest struct {
	// 当前访问令牌。
	Token string `json:"token"`
}
