package oauthrouters

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
	"goauth/middleware"
)

func LoadOAuthIntrospectRoutes(router *gin.Engine, oauthIntrospectController *oauthcontrollers.OAuthIntrospectController, m *middleware.Manager) {
	oauthIntrospectRouter := router.Group("/api/v1/oauth/introspect")
	oauthIntrospectRouter.POST("", m.Auth(), oauthIntrospectController.IntrospectAccessTokenHandler)
}