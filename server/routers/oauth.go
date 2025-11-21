package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
	"goauth/middleware"
	"goauth/services"
)

func LoadOAuthRoutes(router *gin.Engine, ctrl *controllers.OAuthController, authService *services.AuthService) {
	m := middleware.NewManager(authService)
	oauthRouter := router.Group("/api/v1/oauth")
	oauthRouter.GET("/authorization", m.Auth(), ctrl.AuthorizationCodeHandler)
	oauthRouter.POST("/token", ctrl.ExchangeAccessTokenHandler)
}