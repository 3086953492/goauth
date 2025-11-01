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

	AuthService    *services.AuthService
	AuthController *controllers.AuthController
}

func NewContainer() *Container {
	c := &Container{}
	db := database.GetGlobalDB()

	c.UserRepository = repositories.NewUserRepository(db)
	c.UserService = services.NewUserService(c.UserRepository)
	
	c.AuthService = services.NewAuthService(c.UserRepository)
	c.AuthController = controllers.NewAuthController(c.AuthService)

	return c
}
