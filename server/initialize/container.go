package initialize

import (
	"github.com/3086953492/gokit/database"

	"goauth/controllers"
	"goauth/repositories"
	"goauth/services"
)

type Container struct {
	UserRepository *repositories.UserRepository
	UserService    *services.UserService
	UserController *controllers.UserController

	AuthService    *services.AuthService
	AuthController *controllers.AuthController

	OAuthClientRepository *repositories.OAuthClientRepository
	OAuthClientService *services.OAuthClientService
	OAuthClientController *controllers.OAuthClientController
}

func NewContainer() *Container {
	c := &Container{}
	db := database.GetGlobalDB()

	c.UserRepository = repositories.NewUserRepository(db)
	c.UserService = services.NewUserService(c.UserRepository)
	c.UserController = controllers.NewUserController(c.UserService)

	c.AuthService = services.NewAuthService(c.UserRepository)
	c.AuthController = controllers.NewAuthController(c.AuthService)

	c.OAuthClientRepository = repositories.NewOAuthClientRepository(db)
	c.OAuthClientService = services.NewOAuthClientService(c.OAuthClientRepository)
	c.OAuthClientController = controllers.NewOAuthClientController(c.OAuthClientService)

	return c
}
