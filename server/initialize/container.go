package initialize

import (
	"github.com/3086953492/gokit/storage"
	"github.com/3086953492/gokit/validator"
	"gorm.io/gorm"

	"goauth/controllers"
	"goauth/repositories"
	"goauth/services"
	"goauth/validations"
)

type Container struct {
	UserRepository *repositories.UserRepository
	UserService    *services.UserService
	UserController *controllers.UserController
	UserValidator  *validations.UserValidators

	AuthService    *services.AuthService
	AuthController *controllers.AuthController

	OAuthClientRepository *repositories.OAuthClientRepository
	OAuthClientService    *services.OAuthClientService
	OAuthClientController *controllers.OAuthClientController

	OAuthAuthorizationCodeRepository *repositories.OAuthAuthorizationCodeRepository
	OAuthAuthorizationCodeService    *services.OAuthAuthorizationCodeService

	OAuthRefreshTokenRepository *repositories.OAuthRefreshTokenRepository
	OAuthRefreshTokenService    *services.OAuthRefreshTokenService

	OAuthAccessTokenRepository *repositories.OAuthAccessTokenRepository
	OAuthAccessTokenService    *services.OAuthAccessTokenService

	OAuthController *controllers.OAuthController

	ValidatorManager *validator.Manager
}

func NewContainer(db *gorm.DB, storageManager *storage.Manager, validatorManager *validator.Manager) *Container {
	c := &Container{}

	c.UserRepository = repositories.NewUserRepository(db)
	c.UserService = services.NewUserService(c.UserRepository, storageManager)
	c.UserController = controllers.NewUserController(c.UserService, validatorManager)
	c.UserValidator = validations.NewUserValidators(c.UserService)

	c.AuthService = services.NewAuthService(c.UserRepository, c.UserService)
	c.AuthController = controllers.NewAuthController(c.AuthService, validatorManager)

	c.OAuthClientRepository = repositories.NewOAuthClientRepository(db)
	c.OAuthClientService = services.NewOAuthClientService(c.OAuthClientRepository)
	c.OAuthClientController = controllers.NewOAuthClientController(c.OAuthClientService, validatorManager)

	c.OAuthAuthorizationCodeRepository = repositories.NewOAuthAuthorizationCodeRepository(db)
	c.OAuthAuthorizationCodeService = services.NewOAuthAuthorizationCodeService(c.OAuthAuthorizationCodeRepository, c.OAuthClientService)

	c.OAuthRefreshTokenRepository = repositories.NewOAuthRefreshTokenRepository(db)
	c.OAuthRefreshTokenService = services.NewOAuthRefreshTokenService(c.OAuthRefreshTokenRepository)

	c.OAuthAccessTokenRepository = repositories.NewOAuthAccessTokenRepository(db)
	c.OAuthAccessTokenService = services.NewOAuthAccessTokenService(db, c.OAuthAccessTokenRepository, c.OAuthAuthorizationCodeService, c.UserService, c.OAuthClientService, c.OAuthRefreshTokenService)

	c.OAuthController = controllers.NewOAuthController(c.OAuthAuthorizationCodeService, c.OAuthAccessTokenService, c.OAuthClientService)

	c.ValidatorManager = validatorManager

	return c
}
