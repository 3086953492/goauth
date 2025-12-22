package oauthrouters

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers/oauth"
	"goauth/middleware"
)

func LoadOAuthClientRoutes(router *gin.Engine, ctrl *oauthcontrollers.OAuthClientController, m *middleware.Manager) {
	oauthClientRouter := router.Group("/api/v1/oauth/clients")
	oauthClientRouter.POST("", m.Auth(), m.Role("admin"), ctrl.CreateOAuthClientHandler)
	oauthClientRouter.GET("", m.Auth(), m.Role("admin"), ctrl.ListOAuthClientsHandler)
	oauthClientRouter.GET("/:id", m.Auth(), m.Role("admin"), ctrl.GetOAuthClientHandler)
	oauthClientRouter.PATCH("/:id", m.Auth(), m.Role("admin"), ctrl.UpdateOAuthClientHandler)
	oauthClientRouter.DELETE("/:id", m.Auth(), m.Role("admin"), ctrl.DeleteOAuthClientHandler)
}