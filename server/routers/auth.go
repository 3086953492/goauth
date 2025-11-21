package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
	"goauth/middleware"
	"goauth/services"
)

func LoadAuthRoutes(router *gin.Engine,ctrl *controllers.AuthController, authService *services.AuthService) {
	m := middleware.NewManager(authService)
	authRouter := router.Group("/api/v1/auth")
	authRouter.POST("/login", ctrl.LoginHandler)
	authRouter.POST("/logout",m.Auth(), ctrl.LogoutHandler)
}