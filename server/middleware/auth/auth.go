package auth

import (
	"strconv"

	"github.com/3086953492/gokit/cookie"
	"github.com/3086953492/gokit/jwt"
	"github.com/3086953492/gokit/ginx"
	"github.com/gin-gonic/gin"

	"goauth/utils"
)

// Auth 访问令牌验证中间件
func AuthTokenMiddleware(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Cookie 中获取令牌
		token, err := cookie.GetAccessToken(c)
		if err != nil || token == "" {
			ginx.Fail(c, 401, "UNAUTHORIZED", "令牌为空", "about:blank")
			c.Abort()
			return
		}

		claims, err := jwtManager.ParseToken(token)
		if err != nil {
			ginx.Fail(c, 401, "UNAUTHORIZED", "令牌验证失败", "about:blank")
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(claims.UserID, 10, 64)
		if err != nil {
			ginx.Fail(c, 401, "UNAUTHORIZED", "用户ID格式错误", "about:blank")
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
			ginx.Fail(c, 403, "FORBIDDEN", "无权限", "about:blank")
			c.Abort()
			return
		}
		c.Next()
	}
}

func ResourceOwnerMiddleware(source string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsResourceOwner(c, source) {
			ginx.Fail(c, 403, "FORBIDDEN", "无权限", "about:blank")
			c.Abort()
			return
		}
		c.Next()
	}
}
