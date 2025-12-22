package initialize

import (
	"github.com/gin-gonic/gin"

	"goauth/routers"
	"goauth/routers/oauth"
)

func InitRouters(container *Container) *gin.Engine {
	router := gin.Default()

	container.MiddlewareManager.LoadGlobal(router)

	// 注册路由
	routers.LoadAuthRoutes(router, container.AuthController, container.MiddlewareManager)
	routers.LoadUserRoutes(router, container.UserController, container.MiddlewareManager)
	
	oauthrouters.LoadOAuthClientRoutes(router, container.OAuthClientController, container.MiddlewareManager)
	oauthrouters.LoadOAuthAuthorizeRoutes(router, container.OAuthAuthorizeController, container.MiddlewareManager)
	oauthrouters.LoadOAuthIntrospectRoutes(router, container.OAuthIntrospectController, container.MiddlewareManager)
	oauthrouters.LoadOAuthTokenRoutes(router, container.OAuthTokenController, container.MiddlewareManager)

	return router
}
