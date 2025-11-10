package initialize

import (
	"github.com/gin-gonic/gin"

	"goauth/middleware"
	"goauth/routers"
)

func InitRouters(container *Container) *gin.Engine {
	router := gin.Default()

	middlewareManager := middleware.NewManager()
	middlewareManager.LoadGlobal(router)

	// 注册路由
	routers.LoadAuthRoutes(router, container.AuthController)
	routers.LoadUserRoutes(router, container.UserController)
	routers.LoadOAuthClientRoutes(router, container.OAuthClientController)
	
	return router
}
