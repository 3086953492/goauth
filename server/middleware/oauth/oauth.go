package oauth

import (
	"strings"

	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/response"
	"github.com/gin-gonic/gin"
)

// OAuth 访问令牌验证中间件
func OAuthTokenMiddleware(requiredScopes ...string) gin.HandlerFunc {
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

		// token := tokenParts[1]



		// c.Set("user_id", user.ID)


		// if err := utils.EnsureScopes(c, requiredScopes...); err != nil {
		// 	response.Error(c, err)
		// 	c.Abort()
		// 	return
		// }

		c.Next()
	}
}
