package auth

import (
	"strconv"

	"github.com/3086953492/gokit/cookie"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/jwt"
	"github.com/3086953492/gokit/response"
	"github.com/gin-gonic/gin"

	"goauth/services"
	"goauth/utils"
)

// Auth 访问令牌验证中间件
func AuthTokenMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Cookie 中获取令牌
		token, err := cookie.GetAccessToken(c)
		if err != nil || token == "" {
			response.Error(c, errors.Unauthorized().Msg("令牌为空").Build())
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			if errors.IsJwtTokenExpiredError(err) {
				accessToken, accessTokenExpire, err := authService.RefreshToken(c.Request.Context(), token)
				if err != nil {
					response.Error(c, err)
					c.Abort()
					return
				}
				cookie.SetAccessToken(c, accessToken, accessTokenExpire, "", "/")
			} else {
				response.Error(c, errors.Unauthorized().Msg("令牌验证失败").Err(err).Build())
				c.Abort()
				return
			}
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
