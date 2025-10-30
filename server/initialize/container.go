package initialize

import (
	"github.com/3086953492/gokit/database"

	"goauth/controllers"
	"goauth/repositories"
	"goauth/services"
)

type Container struct {
	UserController *controllers.UserController

	UserService *services.UserService

	UserRepository *repositories.UserRepository
}

func NewContainer() *Container {
	c := &Container{}
	db := database.GetGlobalDB()

	c.UserRepository = repositories.NewUserRepository(db)

	c.UserService = services.NewUserService(c.UserRepository)

	c.UserController = controllers.NewUserController(c.UserService)

	return c
}
