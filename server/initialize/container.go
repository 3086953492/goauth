package initialize

import (
	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/jwt"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/redis"
	"github.com/3086953492/gokit/storage"
	"github.com/3086953492/gokit/validator"
	"gorm.io/gorm"

	"goauth/controllers"
	"goauth/controllers/oauth"
	"goauth/middleware"
	"goauth/repositories"
	"goauth/repositories/oauth"
	"goauth/services"
	"goauth/services/oauth"
	"goauth/validations"
)

type Container struct {
	JwtManager *jwt.Manager
	LogManager *logger.Manager

	UserRepository *repositories.UserRepository
	UserService    *services.UserService
	UserController *controllers.UserController
	UserValidator  *validations.UserValidators

	AuthService    *services.AuthService
	AuthController *controllers.AuthController

	OAuthClientRepository *oauthrepositories.OAuthClientRepository
	OAuthClientService    *oauthservices.OAuthClientService
	OAuthClientController *oauthcontrollers.OAuthClientController

	OAuthAuthorizationCodeRepository *oauthrepositories.OAuthAuthorizationCodeRepository
	OAuthAuthorizeService            *oauthservices.OAuthAuthorizeService
	OAuthAuthorizeController         *oauthcontrollers.OAuthAuthorizeController

	OAuthRefreshTokenRepository *oauthrepositories.OAuthRefreshTokenRepository
	OAuthAccessTokenRepository  *oauthrepositories.OAuthAccessTokenRepository
	OAuthTokenService           *oauthservices.OAuthTokenService
	OAuthTokenController        *oauthcontrollers.OAuthTokenController

	OAuthIntrospectService    *oauthservices.OAuthIntrospectService
	OAuthIntrospectController *oauthcontrollers.OAuthIntrospectController

	ValidatorManager *validator.Manager

	MiddlewareManager *middleware.Manager
}

func NewContainer(db *gorm.DB, storageManager *storage.Manager, validatorManager *validator.Manager, redisMgr *redis.Manager, cacheMgr *cache.Manager, jwtMgr *jwt.Manager, logMgr *logger.Manager, cfg *config.Config) *Container {
	c := &Container{}

	c.LogManager = logMgr

	c.UserRepository = repositories.NewUserRepository(db)
	c.UserService = services.NewUserService(c.UserRepository, storageManager, redisMgr, cacheMgr, c.LogManager)
	c.UserController = controllers.NewUserController(c.UserService, validatorManager)
	c.UserValidator = validations.NewUserValidators(c.UserService)

	c.JwtManager = jwtMgr
	c.JwtManager.SetExtraResolver(c.UserService)

	c.AuthService = services.NewAuthService(c.UserRepository, c.UserService, c.LogManager, c.JwtManager, cfg)
	c.AuthController = controllers.NewAuthController(c.AuthService, validatorManager)

	c.OAuthClientRepository = oauthrepositories.NewOAuthClientRepository(db)
	c.OAuthClientService = oauthservices.NewOAuthClientService(c.OAuthClientRepository, cacheMgr, c.LogManager)
	c.OAuthClientController = oauthcontrollers.NewOAuthClientController(c.OAuthClientService, validatorManager)

	c.OAuthAuthorizationCodeRepository = oauthrepositories.NewOAuthAuthorizationCodeRepository(db)
	c.OAuthAuthorizeService = oauthservices.NewOAuthAuthorizeService(c.OAuthAuthorizationCodeRepository, c.OAuthClientService, cfg, c.LogManager)
	c.OAuthAuthorizeController = oauthcontrollers.NewOAuthAuthorizeController(c.OAuthAuthorizeService, c.OAuthClientService, cfg)

	c.OAuthAccessTokenRepository = oauthrepositories.NewOAuthAccessTokenRepository(db)
	c.OAuthRefreshTokenRepository = oauthrepositories.NewOAuthRefreshTokenRepository(db)
	c.OAuthTokenService = oauthservices.NewOAuthTokenService(db, c.OAuthAccessTokenRepository, c.OAuthAuthorizeService, c.UserService, c.OAuthClientService, c.JwtManager, c.LogManager, cfg)
	c.OAuthTokenController = oauthcontrollers.NewOAuthTokenController(c.OAuthTokenService, c.OAuthClientService)

	c.OAuthIntrospectService = oauthservices.NewOAuthIntrospectService(c.OAuthAccessTokenRepository, c.UserService)
	c.OAuthIntrospectController = oauthcontrollers.NewOAuthIntrospectController(c.OAuthIntrospectService, c.OAuthClientService)

	c.ValidatorManager = validatorManager

	c.MiddlewareManager = middleware.NewManager(&cfg.Middleware, c.JwtManager)

	return c
}
