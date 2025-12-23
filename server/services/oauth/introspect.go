package oauthservices

import (
	"context"
	"time"

	"goauth/dto/oauth"
	"goauth/repositories/oauth"
	"goauth/services"
)

type OAuthIntrospectService struct {
	oauthAccessTokenRepository *oauthrepositories.OAuthAccessTokenRepository
	userService                *services.UserService
}

func NewOAuthIntrospectService(oauthAccessTokenRepository *oauthrepositories.OAuthAccessTokenRepository, userService *services.UserService) *OAuthIntrospectService {
	return &OAuthIntrospectService{oauthAccessTokenRepository: oauthAccessTokenRepository, userService: userService}
}

func (s *OAuthIntrospectService) IntrospectAccessToken(ctx context.Context, accessTokenString string) *oauthdto.IntrospectionResponse {
	// 查询访问令牌
	token, err := s.oauthAccessTokenRepository.Get(ctx, map[string]any{"access_token": accessTokenString})
	if err != nil {
		// 令牌不存在或查询错误，统一返回 active=false（RFC 7662 约定）
		return &oauthdto.IntrospectionResponse{Active: false}
	}

	// 检查令牌是否已撤销
	if token.Revoked {
		return &oauthdto.IntrospectionResponse{Active: false}
	}

	// 检查令牌是否已过期
	if token.ExpiresAt.Before(time.Now()) {
		return &oauthdto.IntrospectionResponse{Active: false}
	}

	// 构造有效令牌的响应
	resp := &oauthdto.IntrospectionResponse{
		Active:    true,
		Scope:     token.Scope,
		ClientID:  token.ClientID,
		TokenType: token.TokenType,
		Exp:       token.ExpiresAt.Unix(),
	}

	// 如果存在用户ID，填充 sub 和 username
	if token.UserID != nil {
		// 查询用户名
		user, err := s.userService.GetUser(ctx, map[string]any{"id": *token.UserID})
		if err == nil {
			resp.Username = user.Username
			resp.Sub = user.Subject
		}
	}

	return resp
}
