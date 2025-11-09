package routers

import (
	"github.com/gin-gonic/gin"

	"goauth/controllers"
	"goauth/middleware"
)

func LoadUserRoutes(router *gin.Engine, ctrl *controllers.UserController) {
	m := middleware.NewManager()
	userRouter := router.Group("/api/v1/users")
	userRouter.POST("", ctrl.CreateUserHandler)
	userRouter.GET("/:user_id", m.Auth(), m.ResourceOwner("param"), ctrl.GetUserHandler)
	userRouter.PATCH("/:user_id", m.Auth(), m.ResourceOwner("param"), ctrl.UpdateUserHandler)
	userRouter.GET("", m.Auth(), m.Role("admin"), ctrl.ListUsersHandler)
}
