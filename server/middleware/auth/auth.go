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

// Auth 访问令牌验证中间件
func AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ctx := c.Request.Context()

		// 从 Authorization 头中获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.Unauthorized().Msg("令牌为空").Build())
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			response.Error(c, errors.Unauthorized().Msg("令牌格式不正确").Build())
			c.Abort()
			return
		}

		token := tokenParts[1]

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
