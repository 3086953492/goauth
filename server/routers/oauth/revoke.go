package oauthrouters

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
)

// LoadOAuthRevokeRoutes 注册令牌撤销路由（RFC7009）
func LoadOAuthRevokeRoutes(router *gin.Engine, oauthRevokeController *oauthcontrollers.OAuthRevokeController) {
	oauthRevokeRouter := router.Group("/api/v1/oauth/revoke")
	oauthRevokeRouter.POST("", oauthRevokeController.RevokeTokenHandler)
}

