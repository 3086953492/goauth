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
	"goauth/middleware"
	"goauth/repositories"
	"goauth/repositories/oauth"
	"goauth/services"
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
	OAuthClientService    *services.OAuthClientService
	OAuthClientController *controllers.OAuthClientController

	OAuthAuthorizationCodeRepository *oauthrepositories.OAuthAuthorizationCodeRepository
	OAuthAuthorizationCodeService    *services.OAuthAuthorizationCodeService

	OAuthRefreshTokenRepository *oauthrepositories.OAuthRefreshTokenRepository
	OAuthRefreshTokenService    *services.OAuthRefreshTokenService

	OAuthAccessTokenRepository *oauthrepositories.OAuthAccessTokenRepository
	OAuthAccessTokenService    *services.OAuthAccessTokenService

	OAuthController *controllers.OAuthController

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
	c.OAuthClientService = services.NewOAuthClientService(c.OAuthClientRepository, cacheMgr, c.LogManager)
	c.OAuthClientController = controllers.NewOAuthClientController(c.OAuthClientService, validatorManager)

	c.OAuthAuthorizationCodeRepository = oauthrepositories.NewOAuthAuthorizationCodeRepository(db)
	c.OAuthAuthorizationCodeService = services.NewOAuthAuthorizationCodeService(c.OAuthAuthorizationCodeRepository, c.OAuthClientService, cfg, c.LogManager)

	c.OAuthRefreshTokenRepository = oauthrepositories.NewOAuthRefreshTokenRepository(db)
	c.OAuthRefreshTokenService = services.NewOAuthRefreshTokenService(c.OAuthRefreshTokenRepository, c.JwtManager, cfg, c.LogManager)

	c.OAuthAccessTokenRepository = oauthrepositories.NewOAuthAccessTokenRepository(db)
	c.OAuthAccessTokenService = services.NewOAuthAccessTokenService(db, c.OAuthAccessTokenRepository, c.OAuthAuthorizationCodeService, c.UserService, c.OAuthClientService, c.OAuthRefreshTokenService, c.JwtManager, c.LogManager, cfg)

	c.OAuthController = controllers.NewOAuthController(c.OAuthAuthorizationCodeService, c.OAuthAccessTokenService, c.OAuthClientService, cfg)

	c.ValidatorManager = validatorManager

	c.MiddlewareManager = middleware.NewManager(&cfg.Middleware, c.JwtManager)

	return c
}
