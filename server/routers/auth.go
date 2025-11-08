package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
)

func LoadAuthRoutes(router *gin.Engine,ctrl *controllers.AuthController) {
	authRouter := router.Group("/api/v1/auth")
	authRouter.POST("/login", ctrl.LoginHandler)
}