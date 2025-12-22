package oauthrouters

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
	"goauth/middleware"
)

func LoadOAuthAuthorizeRoutes(router *gin.Engine, oauthAuthorizeController *oauthcontrollers.OAuthAuthorizeController, m *middleware.Manager) {
	oauthAuthorizeRouter := router.Group("/api/v1/oauth/authorize")
	oauthAuthorizeRouter.GET("", m.Auth(), oauthAuthorizeController.AuthorizationCodeHandler)
}