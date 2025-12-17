package initialize

import (
	"github.com/gin-gonic/gin"

	"goauth/routers"
)

func InitRouters(container *Container) *gin.Engine {
	router := gin.Default()

	container.MiddlewareManager.LoadGlobal(router)

	// 注册路由
	routers.LoadAuthRoutes(router, container.AuthController, container.MiddlewareManager)
	routers.LoadUserRoutes(router, container.UserController, container.MiddlewareManager)
	routers.LoadOAuthClientRoutes(router, container.OAuthClientController, container.MiddlewareManager)
	routers.LoadOAuthRoutes(router, container.OAuthController, container.MiddlewareManager)

	return router
}
