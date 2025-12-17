package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
	"goauth/middleware"
)

func LoadOAuthClientRoutes(router *gin.Engine, ctrl *controllers.OAuthClientController, m *middleware.Manager) {
	oauthClientRouter := router.Group("/api/v1/oauth_clients")
	oauthClientRouter.POST("", m.Auth(), m.Role("admin"), ctrl.CreateOAuthClientHandler)
	oauthClientRouter.GET("", m.Auth(), m.Role("admin"), ctrl.ListOAuthClientsHandler)
	oauthClientRouter.GET("/:id", m.Auth(), m.Role("admin"), ctrl.GetOAuthClientHandler)
	oauthClientRouter.PATCH("/:id", m.Auth(), m.Role("admin"), ctrl.UpdateOAuthClientHandler)
	oauthClientRouter.DELETE("/:id", m.Auth(), m.Role("admin"), ctrl.DeleteOAuthClientHandler)
}