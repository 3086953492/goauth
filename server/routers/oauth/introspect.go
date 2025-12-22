package oauthrouters

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
	"goauth/middleware"
)

func LoadOAuthIntrospectRoutes(router *gin.Engine, oauthTokenController *oauthcontrollers.OAuthTokenController, m *middleware.Manager) {
	oauthIntrospectRouter := router.Group("/api/v1/oauth/introspect")
	oauthIntrospectRouter.POST("", m.Auth(), oauthTokenController.IntrospectAccessTokenHandler)
}