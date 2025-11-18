package services

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/errors"

	"goauth/models"
	"goauth/repositories"
)

type OAuthAuthorizationCodeService struct {
	oauthAuthorizationCodeRepository *repositories.OAuthAuthorizationCodeRepository
	oauthClientService               *OAuthClientService
}

func NewOAuthAuthorizationCodeService(oauthAuthorizationCodeRepository *repositories.OAuthAuthorizationCodeRepository, oauthClientService *OAuthClientService) *OAuthAuthorizationCodeService {
	return &OAuthAuthorizationCodeService{oauthAuthorizationCodeRepository: oauthAuthorizationCodeRepository, oauthClientService: oauthClientService}
}

func (s *OAuthAuthorizationCodeService) GenerateAuthorizationCode(ctx context.Context, userID uint, clientID string, redirectURI string, scope string, expiresIn float64) (string, error) {

	oauthClient, err := s.oauthClientService.GetOAuthClient(ctx, map[string]any{"id": clientID, "status": 1})
	if err != nil {
		return "", err
	}

	if !isRedirectURIValid(redirectURI, oauthClient.RedirectURIs) {
		return "", errors.InvalidInput().Msg("redirect_uri不在客户端的回调地址列表中").Build()
	}

	if !isScopeValid(scope, oauthClient.Scopes) {
		return "", errors.InvalidInput().Msg("scope不在客户端的权限范围列表中").Build()
	}

	codeString, err := crypto.GenerateAuthorizationCode(32)
	if err != nil {
		return "", errors.Internal().Msg("生成授权码失败").Err(err).Build()
	}

	code := &models.OAuthAuthorizationCode{
		Code:        codeString,
		UserID:      userID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
	if err := s.oauthAuthorizationCodeRepository.Create(ctx, code); err != nil {
		return "", errors.Database().Msg("创建OAuth授权码失败").Err(err).Field("code", code).Log()
	}

	return codeString, nil
}

// isRedirectURIValid 验证 redirect_uri 是否在白名单中（精确匹配）
func isRedirectURIValid(redirectURI string, registeredURIsJSON []byte) bool {
	var registeredURIs []string
	if err := json.Unmarshal(registeredURIsJSON, &registeredURIs); err != nil {
		return false
	}

	for _, uri := range registeredURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}

// isScopeValid 验证请求的 scope 是否都在允许的 scope 列表中
func isScopeValid(requestedScope string, allowedScopesJSON []byte) bool {
	if requestedScope == "" {
		return true // 空 scope 视为合法
	}

	var allowedScopes []string
	if err := json.Unmarshal(allowedScopesJSON, &allowedScopes); err != nil {
		return false
	}

	// 将允许的 scope 转为 map 便于查找
	allowedScopeMap := make(map[string]bool)
	for _, scope := range allowedScopes {
		allowedScopeMap[scope] = true
	}

	// 检查请求的每个 scope 是否都在允许列表中
	requestedScopes := strings.Split(requestedScope, " ")
	for _, scope := range requestedScopes {
		scope = strings.TrimSpace(scope)
		if scope != "" && !allowedScopeMap[scope] {
			return false // 有任何一个 scope 不在允许列表中就返回 false
		}
	}

	return true
}

func (s *OAuthAuthorizationCodeService) GetOAuthAuthorizationCode(ctx context.Context, conds map[string]any) (*models.OAuthAuthorizationCode, error) {
	oauthAuthorizationCode, err := s.oauthAuthorizationCodeRepository.Get(ctx, conds)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, err
		}
		return nil, errors.Database().Msg("系统繁忙，请稍后再试").Err(err).Log()
	}
	return oauthAuthorizationCode, nil
}

func (s *OAuthAuthorizationCodeService) MarkAsUsed(ctx context.Context, id uint) error {
	if err := s.oauthAuthorizationCodeRepository.MarkAsUsed(ctx, id); err != nil {
		return errors.Database().Msg("标记授权码为已使用失败").Err(err).Log()
	}
	return nil
}
