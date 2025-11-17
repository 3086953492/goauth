package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
	"goauth/middleware"
)

func LoadOAuthRoutes(router *gin.Engine, ctrl *controllers.OAuthController) {
	m := middleware.NewManager()
	oauthRouter := router.Group("/api/v1/oauth")
	oauthRouter.GET("/authorization", m.Auth(), ctrl.AuthorizationCodeHandler)
}