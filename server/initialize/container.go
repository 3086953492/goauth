package initialize

import (
	"github.com/3086953492/gokit/database"

	"goauth/controllers"
	"goauth/repositories"
	"goauth/services"
)

type Container struct {
	UserRepository *repositories.UserRepository

	AuthController *controllers.AuthController
	AuthService    *services.AuthService
}

func NewContainer() *Container {
	c := &Container{}
	db := database.GetGlobalDB()

	c.UserRepository = repositories.NewUserRepository(db)

	c.AuthController = controllers.NewAuthController(c.AuthService)
	c.AuthService = services.NewAuthService(c.UserRepository)

	return c
}
