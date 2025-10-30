package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
)

func LoadUserRoutes(router *gin.Engine,ctrl *controllers.UserController) {
	userRouter := router.Group("/users")
	userRouter.POST("/create", ctrl.CreateUserHandler)
}