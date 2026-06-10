package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	apperrors "github.com/tohru-shippo/ai-mate-client-gateway/internal/errors"
	"github.com/tohru-shippo/ai-mate-client-gateway/internal/response"
)

const userPrincipalKey = "userPrincipal"

// 已通过鉴权的用户身份。
type UserPrincipal struct {
	// 用户 ID。
	UserID string
	// 用户名。
	Username string
	// 用户角色。
	Role string
	// 用户首选语言。
	Language string
}

// 用户 token 鉴权能力。
type Authenticator interface {
	// 校验访问令牌并返回用户身份。
	Authenticate(ctx context.Context, token string) (*UserPrincipal, error)
}

// 校验 Bearer token 并把用户身份写入 Gin 上下文。
func Auth(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := parseBearerToken(c.GetHeader("Authorization"))
		if token == "" {
			response.FromError(c, apperrors.Unauthorized())
			c.Abort()
			return
		}
		principal, err := authenticator.Authenticate(c.Request.Context(), token)
		if err != nil {
			response.FromError(c, err)
			c.Abort()
			return
		}
		c.Set(userPrincipalKey, principal)
		c.Next()
	}
}

// 从 Gin 上下文读取已鉴权的用户身份。
func MustUser(c *gin.Context) (*UserPrincipal, error) {
	value, exists := c.Get(userPrincipalKey)
	if !exists {
		return nil, errors.New("user principal not found")
	}
	principal, ok := value.(*UserPrincipal)
	if !ok || principal == nil {
		return nil, errors.New("invalid user principal")
	}
	return principal, nil
}

func parseBearerToken(value string) string {
	const prefix = "Bearer "
	value = strings.TrimSpace(value)
	if !strings.HasPrefix(value, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(value, prefix))
}
