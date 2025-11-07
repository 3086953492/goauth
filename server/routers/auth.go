package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
)

func LoadAuthRoutes(router *gin.Engine,ctrl *controllers.AuthController) {
	authRouter := router.Group("/auth")
	authRouter.POST("/register", ctrl.RegisterHandler)
	authRouter.POST("/login", ctrl.LoginHandler)
}