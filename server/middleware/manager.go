package middleware

import (
	"github.com/3086953492/gokit/config/types"
	"github.com/3086953492/gokit/ginx/cookie"
	"github.com/3086953492/gokit/jwt"
	"github.com/gin-gonic/gin"

	"goauth/middleware/auth"
	"goauth/middleware/security"
	"goauth/repositories/oauth"
)

// 中间件管理器
type Manager struct {
	config          *types.MiddlewareConfig
	jwtManager      *jwt.Manager
	cookieMgr       *cookie.TokenCookies
	accessTokenRepo *oauthrepositories.OAuthAccessTokenRepository
}

// 创建管理器（通过注入配置）
func NewManager(
	cfg *types.MiddlewareConfig,
	jwtManager *jwt.Manager,
	cookieMgr *cookie.TokenCookies,
	accessTokenRepo *oauthrepositories.OAuthAccessTokenRepository,
) *Manager {
	return &Manager{
		config:          cfg,
		jwtManager:      jwtManager,
		cookieMgr:       cookieMgr,
		accessTokenRepo: accessTokenRepo,
	}
}

// 加载所有全局中间件
func (m *Manager) LoadGlobal(engine *gin.Engine) {
	// CORS中间件
	engine.Use(m.CORS())

	// 这里以后可以加更多全局中间件
	// engine.Use(m.Logger())
}

// CORS中间件
func (m *Manager) CORS() gin.HandlerFunc {
	return security.NewCORSMiddleware(m.config.CORS)
}

// Auth 默认认证中间件，仅支持 Cookie(JWT)
// 不接受 Bearer token，适用于内部接口
func (m *Manager) Auth() gin.HandlerFunc {
	return auth.AuthCookieMiddleware(m.jwtManager, m.cookieMgr)
}

// AuthBearerOrCookie 支持 Bearer(OAuth) 和 Cookie(JWT) 的认证中间件
// 优先 Bearer，无 Bearer 时回退 Cookie
// 通过 opts 配置允许的 Bearer 主体类型，例如：
//   - auth.BearerAllowAll()   允许所有 Bearer 类型
//   - auth.BearerAllowUser()  仅允许 Bearer-User
//   - auth.BearerAllowClient() 仅允许 Bearer-Client（client_credentials）
func (m *Manager) AuthBearerOrCookie(opts ...auth.BearerOption) gin.HandlerFunc {
	return auth.AuthBearerOrCookieMiddleware(m.jwtManager, m.cookieMgr, m.accessTokenRepo, opts...)
}

func (m *Manager) Role(requiredRole string) gin.HandlerFunc {
	return auth.RoleMiddleware(requiredRole)
}

func (m *Manager) ResourceOwner(source string) gin.HandlerFunc {
	return auth.ResourceOwnerMiddleware(source)
}

// Scope 检查 OAuth scope 是否包含所需权限
func (m *Manager) Scope(requiredScopes ...string) gin.HandlerFunc {
	return auth.ScopeMiddleware(requiredScopes...)
}

// RequireUser 要求当前主体必须是用户（非客户端）
func (m *Manager) RequireUser() gin.HandlerFunc {
	return auth.RequireUserMiddleware()
}
