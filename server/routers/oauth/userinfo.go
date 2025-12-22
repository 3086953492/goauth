package oauthrouters

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
)

func LoadOAuthUserInfoRoutes(router *gin.Engine, oauthUserInfoController *oauthcontrollers.OAuthUserInfoController) {
	oauthUserInfoRouter := router.Group("/api/v1/oauth/userinfo")
	oauthUserInfoRouter.GET("", oauthUserInfoController.GetUserInfoHandler)
}
