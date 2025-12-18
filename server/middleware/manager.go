package middleware

import (
	"github.com/3086953492/gokit/config/types"
	"github.com/3086953492/gokit/jwt"
	"github.com/gin-gonic/gin"

	"goauth/middleware/auth"
	"goauth/middleware/security"
)

// 中间件管理器
type Manager struct {
	config *types.MiddlewareConfig
	jwtManager *jwt.Manager
}

// 创建管理器（通过注入配置）
func NewManager(cfg *types.MiddlewareConfig, jwtManager *jwt.Manager) *Manager {
	return &Manager{
		config: cfg,
		jwtManager: jwtManager,
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

func (m *Manager) Auth() gin.HandlerFunc {
	return auth.AuthTokenMiddleware(m.jwtManager)
}

func (m *Manager) Role(requiredRole string) gin.HandlerFunc {
	return auth.RoleMiddleware(requiredRole)
}

func (m *Manager) ResourceOwner(source string) gin.HandlerFunc {
	return auth.ResourceOwnerMiddleware(source)
}
