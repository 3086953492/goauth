package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
	"goauth/middleware"
)

func LoadOAuthRoutes(router *gin.Engine, oauthAuthorizationController *oauthcontrollers.OAuthAuthorizationController, oauthTokenController *oauthcontrollers.OAuthTokenController, m *middleware.Manager) {
	oauthRouter := router.Group("/api/v1/oauth")
	oauthRouter.GET("/authorization", m.Auth(), oauthAuthorizationController.AuthorizationCodeHandler)
	oauthRouter.POST("/token", oauthTokenController.ExchangeAccessTokenHandler)
	oauthRouter.POST("/introspect", oauthTokenController.IntrospectAccessTokenHandler)
}