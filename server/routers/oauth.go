package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
)

func LoadOAuthRoutes(router *gin.Engine, ctrl *controllers.OAuthController) {
	oauthRouter := router.Group("/api/v1/oauth")
	oauthRouter.POST("/authorization", ctrl.AuthorizationCodeHandler)
}