package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
	"goauth/middleware"
)

func LoadAuthRoutes(router *gin.Engine, ctrl *controllers.AuthController, m *middleware.Manager) {
	authRouter := router.Group("/api/v1/auth")
	authRouter.POST("/login", ctrl.LoginHandler)
	authRouter.POST("/logout", m.Auth(), ctrl.LogoutHandler)
	authRouter.POST("/refresh_token", ctrl.RefreshTokenHandler)
}
