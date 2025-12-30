package auth

import (
	"strconv"
	"strings"
	"time"

	"github.com/3086953492/gokit/ginx/cookie"
	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/jwt"
	"github.com/gin-gonic/gin"

	"goauth/repositories/oauth"
	"goauth/utils"
)

// PrincipalKind 表示认证主体类型
type PrincipalKind string

const (
	PrincipalKindUser   PrincipalKind = "user"
	PrincipalKindClient PrincipalKind = "client"
)

// Principal 统一认证主体，写入 gin.Context
type Principal struct {
	Kind     PrincipalKind // user 或 client
	UserID   uint64        // 用户ID（仅 user 主体有效）
	ClientID string        // OAuth 客户端ID（仅 bearer 认证有效）
	Scope    string        // OAuth scope（仅 bearer 认证有效）
	Role     string        // 角色（仅 cookie 认证有效）
}

// setUserPrincipal 设置用户主体到 context
func setUserPrincipal(c *gin.Context, userID uint64, role string) {
	c.Set("principal_kind", string(PrincipalKindUser))
	c.Set("user_id", userID)
	c.Set("role", role)
	c.Set("client_id", "")
	c.Set("scope", "")
}

// setClientPrincipal 设置客户端主体到 context（client_credentials 模式，无用户）
func setClientPrincipal(c *gin.Context, clientID, scope string) {
	c.Set("principal_kind", string(PrincipalKindClient))
	c.Set("user_id", uint64(0))
	c.Set("role", "")
	c.Set("client_id", clientID)
	c.Set("scope", scope)
}

// setBearerUserPrincipal 设置 bearer token 的用户主体（授权码/刷新令牌模式）
func setBearerUserPrincipal(c *gin.Context, userID uint64, clientID, scope string) {
	c.Set("principal_kind", string(PrincipalKindUser))
	c.Set("user_id", userID)
	c.Set("role", "")
	c.Set("client_id", clientID)
	c.Set("scope", scope)
}

// ============================================================================
// Bearer 策略选项
// ============================================================================

// bearerPolicy 控制允许的 Bearer 主体类型
type bearerPolicy struct {
	allowUser   bool
	allowClient bool
}

// BearerOption 配置 Bearer 策略
type BearerOption func(*bearerPolicy)

// BearerAllowUser 允许 Bearer-User（授权码/刷新令牌模式）
func BearerAllowUser() BearerOption {
	return func(p *bearerPolicy) {
		p.allowUser = true
	}
}

// BearerAllowClient 允许 Bearer-Client（client_credentials 模式）
func BearerAllowClient() BearerOption {
	return func(p *bearerPolicy) {
		p.allowClient = true
	}
}

// BearerAllowAll 允许所有 Bearer 类型（User + Client）
func BearerAllowAll() BearerOption {
	return func(p *bearerPolicy) {
		p.allowUser = true
		p.allowClient = true
	}
}

// ============================================================================
// Cookie-Only 中间件（默认）
// ============================================================================

// AuthCookieMiddleware 仅 Cookie(JWT) 认证中间件
// 不接受 Bearer token，仅校验内部 Cookie 令牌
func AuthCookieMiddleware(
	jwtManager *jwt.Manager,
	cookieMgr *cookie.TokenCookies,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authenticateByCookie(c, jwtManager, cookieMgr) {
			c.Next()
			return
		}
		// 认证失败，已在 authenticateByCookie 中返回 401
	}
}

// ============================================================================
// BearerOrCookie 中间件（按组显式启用）
// ============================================================================

// AuthBearerOrCookieMiddleware 支持 Bearer(OAuth) 和 Cookie(JWT) 的认证中间件
// 优先 Bearer，无 Bearer 时回退 Cookie
// 通过 opts 配置允许的 Bearer 主体类型
func AuthBearerOrCookieMiddleware(
	jwtManager *jwt.Manager,
	cookieMgr *cookie.TokenCookies,
	accessTokenRepo *oauthrepositories.OAuthAccessTokenRepository,
	opts ...BearerOption,
) gin.HandlerFunc {
	// 构建策略（默认不允许任何 Bearer）
	policy := &bearerPolicy{}
	for _, opt := range opts {
		opt(policy)
	}

	return func(c *gin.Context) {
		// 1. 优先尝试 Bearer 认证
		authHeader := c.GetHeader("Authorization")
		if token, found := strings.CutPrefix(authHeader, "Bearer "); found && token != "" {
			if authenticateByBearerWithPolicy(c, token, accessTokenRepo, policy) {
				c.Next()
				return
			}
			// Bearer 认证失败，已在 authenticateByBearerWithPolicy 中返回 401/403
			return
		}

		// 2. 回退到 Cookie 认证
		if authenticateByCookie(c, jwtManager, cookieMgr) {
			c.Next()
			return
		}
		// Cookie 认证失败，已在 authenticateByCookie 中返回 401
	}
}

// authenticateByBearerWithPolicy 通过 Bearer token 认证（查库判活），并根据策略过滤主体类型
// 返回 true 表示认证成功，false 表示失败（已返回 401/403）
func authenticateByBearerWithPolicy(
	c *gin.Context,
	token string,
	accessTokenRepo *oauthrepositories.OAuthAccessTokenRepository,
	policy *bearerPolicy,
) bool {
	// 查询 access token
	accessToken, err := accessTokenRepo.Get(c.Request.Context(), map[string]any{"access_token": token})
	if err != nil {
		problem.Fail(c, 401, "UNAUTHORIZED", "令牌无效", "about:blank")
		c.Abort()
		return false
	}

	// 检查是否已撤销
	if accessToken.Revoked {
		problem.Fail(c, 401, "UNAUTHORIZED", "令牌已撤销", "about:blank")
		c.Abort()
		return false
	}

	// 检查是否已过期
	if accessToken.ExpiresAt.Before(time.Now()) {
		problem.Fail(c, 401, "UNAUTHORIZED", "令牌已过期", "about:blank")
		c.Abort()
		return false
	}

	// 根据 UserID 判断是用户主体还是客户端主体
	isUserToken := accessToken.UserID != nil && *accessToken.UserID > 0

	// 按策略过滤
	if isUserToken && !policy.allowUser {
		problem.Fail(c, 403, "FORBIDDEN", "此接口不允许用户令牌访问", "about:blank")
		c.Abort()
		return false
	}
	if !isUserToken && !policy.allowClient {
		problem.Fail(c, 403, "FORBIDDEN", "此接口不允许客户端凭证访问", "about:blank")
		c.Abort()
		return false
	}

	// 设置主体
	if isUserToken {
		setBearerUserPrincipal(c, uint64(*accessToken.UserID), accessToken.ClientID, accessToken.Scope)
	} else {
		setClientPrincipal(c, accessToken.ClientID, accessToken.Scope)
	}

	return true
}

// authenticateByCookie 通过 Cookie 中的 JWT 认证
// 返回 true 表示认证成功，false 表示失败（已返回 401）
func authenticateByCookie(
	c *gin.Context,
	jwtManager *jwt.Manager,
	cookieMgr *cookie.TokenCookies,
) bool {
	// 从 Cookie 中获取令牌
	token, err := cookieMgr.GetAccess(c)
	if err != nil || token == "" {
		problem.Fail(c, 401, "UNAUTHORIZED", "令牌为空", "about:blank")
		c.Abort()
		return false
	}

	claims, err := jwtManager.ParseToken(token)
	if err != nil {
		problem.Fail(c, 401, "UNAUTHORIZED", "令牌验证失败", "about:blank")
		c.Abort()
		return false
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		problem.Fail(c, 401, "UNAUTHORIZED", "用户ID格式错误", "about:blank")
		c.Abort()
		return false
	}

	// 获取角色
	role := ""
	if r, ok := claims.Extra["role"]; ok {
		if roleStr, ok := r.(string); ok {
			role = roleStr
		}
	}

	setUserPrincipal(c, userID, role)
	return true
}

// ============================================================================
// 鉴权中间件
// ============================================================================

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsRole(c, requiredRole) {
			problem.Fail(c, 403, "FORBIDDEN", "无权限", "about:blank")
			c.Abort()
			return
		}
		c.Next()
	}
}

func ResourceOwnerMiddleware(source string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsResourceOwner(c, source) {
			problem.Fail(c, 403, "FORBIDDEN", "无权限", "about:blank")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ScopeMiddleware 检查 OAuth scope 是否包含所需权限
func ScopeMiddleware(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopeVal, exists := c.Get("scope")
		if !exists {
			// 没有 scope（可能是 cookie 认证），默认通过
			c.Next()
			return
		}

		scope, ok := scopeVal.(string)
		if !ok || scope == "" {
			// scope 为空，检查是否需要 scope
			if len(requiredScopes) > 0 {
				problem.Fail(c, 403, "FORBIDDEN", "缺少必要的权限范围", "about:blank")
				c.Abort()
				return
			}
			c.Next()
			return
		}

		// 解析当前 scope
		currentScopes := strings.Split(scope, " ")
		scopeSet := make(map[string]bool)
		for _, s := range currentScopes {
			s = strings.TrimSpace(s)
			if s != "" {
				scopeSet[s] = true
			}
		}

		// 检查是否包含所有必需的 scope
		for _, required := range requiredScopes {
			if !scopeSet[required] {
				problem.Fail(c, 403, "FORBIDDEN", "缺少必要的权限范围: "+required, "about:blank")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// RequireUserMiddleware 要求当前主体必须是用户（非客户端）
func RequireUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		kindVal, exists := c.Get("principal_kind")
		if !exists {
			problem.Fail(c, 403, "FORBIDDEN", "无法确定认证主体类型", "about:blank")
			c.Abort()
			return
		}

		kind, ok := kindVal.(string)
		if !ok || kind != string(PrincipalKindUser) {
			problem.Fail(c, 403, "FORBIDDEN", "此接口仅限用户访问", "about:blank")
			c.Abort()
			return
		}

		// 确保 user_id 有效
		userIDVal, exists := c.Get("user_id")
		if !exists {
			problem.Fail(c, 403, "FORBIDDEN", "用户ID不存在", "about:blank")
			c.Abort()
			return
		}

		userID, ok := userIDVal.(uint64)
		if !ok || userID == 0 {
			problem.Fail(c, 403, "FORBIDDEN", "用户ID无效", "about:blank")
			c.Abort()
			return
		}

		c.Next()
	}
}
