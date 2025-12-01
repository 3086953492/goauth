package services

import (
	"context"
	"strconv"
	"time"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/jwt"

	"goauth/dto"
	"goauth/models"
	"goauth/repositories"
	"goauth/utils"
)

type OAuthAccessTokenService struct {
	oauthAccessTokenRepository *repositories.OAuthAccessTokenRepository

	oauthAuthorizationCodeService *OAuthAuthorizationCodeService

	userService *UserService

	oauthClientService *OAuthClientService

	oauthRefreshTokenService *OAuthRefreshTokenService
}

func NewOAuthAccessTokenService(oauthAccessTokenRepository *repositories.OAuthAccessTokenRepository, oauthAuthorizationCodeService *OAuthAuthorizationCodeService, userService *UserService, oauthClientService *OAuthClientService, oauthRefreshTokenService *OAuthRefreshTokenService) *OAuthAccessTokenService {
	return &OAuthAccessTokenService{oauthAccessTokenRepository: oauthAccessTokenRepository, oauthAuthorizationCodeService: oauthAuthorizationCodeService, userService: userService, oauthClientService: oauthClientService, oauthRefreshTokenService: oauthRefreshTokenService}
}

func (s *OAuthAccessTokenService) ExchangeAccessToken(ctx context.Context, form *dto.ExchangeAccessTokenForm, clientID, clientSecret string) (*dto.ExchangeAccessTokenResponse, error) {

	oauthClient, err := s.oauthClientService.GetOAuthClient(ctx, map[string]any{"id": clientID, "client_secret": clientSecret})
	if err != nil {
		return nil, err
	}

	if form.GrantType != "authorization_code" || !utils.IsGrantTypeValid("authorization_code", oauthClient.GrantTypes) {
		return nil, errors.InvalidInput().Msg("授权类型不支持").Build()
	}

	oauthAuthorizationCode, err := s.oauthAuthorizationCodeService.GetOAuthAuthorizationCode(ctx, map[string]any{"code": form.Code})
	if err != nil {
		return nil, err
	}

	if oauthAuthorizationCode.RedirectURI != form.RedirectURI {
		return nil, errors.InvalidInput().Msg("授权码回调地址不匹配").Build()
	}

	if oauthAuthorizationCode.Used {
		return nil, errors.InvalidInput().Msg("授权码已使用").Build()
	}

	if oauthAuthorizationCode.ExpiresAt.Before(time.Now()) {
		return nil, errors.InvalidInput().Msg("授权码已过期").Build()
	}

	if oauthAuthorizationCode.ClientID != clientID {
		return nil, errors.InvalidInput().Msg("授权码客户端ID不匹配").Build()
	}

	if err := s.oauthAuthorizationCodeService.MarkAsUsed(ctx, oauthAuthorizationCode.ID); err != nil {
		return nil, err
	}

	user, err := s.userService.GetUser(ctx, map[string]any{"id": oauthAuthorizationCode.UserID})
	if err != nil {
		return nil, err
	}

	accessTokenString, err := jwt.GenerateToken(strconv.FormatUint(uint64(user.ID), 10), user.Username, map[string]any{"role": user.Role})
	if err != nil {
		return nil, errors.Internal().Msg("生成访问令牌失败").Err(err).Log()
	}

	accessToken := &models.OAuthAccessToken{
		AccessToken: accessTokenString,
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(config.GetGlobalConfig().OAuth.AccessTokenExpire),
		ClientID:    oauthAuthorizationCode.ClientID,
		Scope:       oauthAuthorizationCode.Scope,
		UserID:      &oauthAuthorizationCode.UserID,
	}

	refreshToken, err := s.oauthRefreshTokenService.GenerateRefreshToken(ctx, accessToken.ID, oauthAuthorizationCode.ClientID, oauthAuthorizationCode.Scope, oauthAuthorizationCode.UserID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	if err := s.oauthAccessTokenRepository.Create(ctx, accessToken); err != nil {
		return nil, errors.Database().Msg("创建OAuth访问令牌失败").Err(err).Field("accessToken", accessToken).Log()
	}
	return &dto.ExchangeAccessTokenResponse{
		AccessToken: dto.OAuthAccessTokenResponse{
			AccessToken: accessTokenString,
			ExpiresIn:   int(config.GetGlobalConfig().OAuth.AccessTokenExpire.Seconds()),
		},
		RefreshToken: dto.OAuthRefreshTokenResponse{
			RefreshToken: refreshToken,
			ExpiresIn:    int(config.GetGlobalConfig().OAuth.RefreshTokenExpire.Seconds()),
		},
		Scope:     accessToken.Scope,
		TokenType: "Bearer",
	}, nil
}
