package oauthrouters

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
	"goauth/middleware"
)

func LoadOAuthTokenRoutes(router *gin.Engine, oauthTokenController *oauthcontrollers.OAuthTokenController, m *middleware.Manager) {
	oauthTokenRouter := router.Group("/api/v1/oauth/token")
	oauthTokenRouter.POST("", oauthTokenController.ExchangeAccessTokenHandler)
}
