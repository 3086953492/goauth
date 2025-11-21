package auth

import (
	"strconv"
	"strings"

	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/jwt"
	"github.com/3086953492/gokit/response"
	"github.com/gin-gonic/gin"

	"goauth/utils"
)

// extractToken 从请求头或 Cookie 中提取访问令牌
// 优先级：Authorization 头 > Authorization Cookie > access_token Cookie
func extractToken(c *gin.Context) string {
	// 1. 尝试从 Authorization 请求头获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
			return tokenParts[1]
		}
	}

	// 2. 尝试从 Authorization Cookie 获取
	authCookie, err := c.Cookie("Authorization")
	if err == nil && authCookie != "" {
		tokenParts := strings.Split(authCookie, " ")
		if len(tokenParts) == 2 && strings.ToLower(tokenParts[0]) == "bearer" {
			return tokenParts[1]
		}
	}

	// 3. 尝试从 access_token Cookie 获取（纯 token，无 Bearer 前缀）
	accessToken, err := c.Cookie("access_token")
	if err == nil && accessToken != "" {
		return accessToken
	}

	return ""
}

// Auth 访问令牌验证中间件
func AuthTokenMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ctx := c.Request.Context()

		// 从请求头或 Cookie 中获取令牌
		token := extractToken(c)
		if token == "" {
			response.Error(c, errors.Unauthorized().Msg("令牌为空").Build())
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Error(c, errors.Unauthorized().Msg("令牌验证失败").Build())
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			response.Error(c, errors.Unauthorized().Msg("用户ID格式错误").Build())
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("role", claims.Extra["role"])

		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsRole(c, requiredRole) {
			response.Error(c, errors.Forbidden().Msg("无权限").Build())
			c.Abort()
			return
		}
		c.Next()
	}
}

func ResourceOwnerMiddleware(source string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsResourceOwner(c, source) {
			response.Error(c, errors.Forbidden().Msg("无权限").Build())
			c.Abort()
			return
		}
		c.Next()
	}
}
