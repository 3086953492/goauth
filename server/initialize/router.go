package initialize

import (
	"github.com/gin-gonic/gin"

	"goauth/middleware"
	"goauth/routers"
)

func InitRouters(container *Container) *gin.Engine {
	router := gin.Default()

	middlewareManager := middleware.NewManager(container.AuthService)
	middlewareManager.LoadGlobal(router)

	// 注册路由
	routers.LoadAuthRoutes(router, container.AuthController)
	routers.LoadUserRoutes(router, container.UserController, container.AuthService)
	routers.LoadOAuthClientRoutes(router, container.OAuthClientController, container.AuthService)
	routers.LoadOAuthRoutes(router, container.OAuthController, container.AuthService)
	
	return router
}
