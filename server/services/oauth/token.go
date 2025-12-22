package oauthservices

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/jwt"
	"github.com/3086953492/gokit/logger"
	"gorm.io/gorm"

	"goauth/dto/oauth"
	"goauth/models/oauth"
	"goauth/repositories/oauth"
	"goauth/services"
	"goauth/utils"
)

type OAuthTokenService struct {
	db *gorm.DB

	oauthAccessTokenRepository  *oauthrepositories.OAuthAccessTokenRepository
	oauthRefreshTokenRepository *oauthrepositories.OAuthRefreshTokenRepository

	oauthAuthorizeService *OAuthAuthorizeService

	userService *services.UserService

	oauthClientService *OAuthClientService

	jwtManager *jwt.Manager
	logMgr     *logger.Manager
	cfg        *config.Config
}

func NewOAuthTokenService(
	db *gorm.DB,
	oauthAccessTokenRepository *oauthrepositories.OAuthAccessTokenRepository,
	oauthAuthorizeService *OAuthAuthorizeService,
	userService *services.UserService,
	oauthClientService *OAuthClientService,
	jwtManager *jwt.Manager,
	logMgr *logger.Manager,
	cfg *config.Config,
) *OAuthTokenService {
	return &OAuthTokenService{
		db:                            db,
		oauthAccessTokenRepository:    oauthAccessTokenRepository,
		oauthAuthorizeService:         oauthAuthorizeService,
		userService:                   userService,
		oauthClientService:            oauthClientService,
		jwtManager:                    jwtManager,
		logMgr:                        logMgr,
		cfg:                           cfg,
	}
}

func (s *OAuthTokenService) ExchangeAccessToken(ctx context.Context, form *oauthdto.ExchangeAccessTokenForm, clientID, clientSecret string) (*oauthdto.ExchangeAccessTokenResponse, error) {

	oauthClient, err := s.oauthClientService.GetOAuthClient(ctx, map[string]any{"id": clientID, "client_secret": clientSecret})
	if err != nil {
		return nil, err
	}

	if form.GrantType != "authorization_code" || !utils.IsGrantTypeValid("authorization_code", oauthClient.GrantTypes) {
		return nil, errors.New("授权类型不支持")
	}

	oauthAuthorizationCode, err := s.oauthAuthorizeService.GetOAuthAuthorizationCode(ctx, map[string]any{"code": form.Code})
	if err != nil {
		return nil, err
	}

	if oauthAuthorizationCode.RedirectURI != form.RedirectURI {
		return nil, errors.New("授权码回调地址不匹配")
	}

	if oauthAuthorizationCode.Used {
		return nil, errors.New("授权码已使用")
	}

	if oauthAuthorizationCode.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("授权码已过期")
	}

	if oauthAuthorizationCode.ClientID != clientID {
		return nil, errors.New("授权码客户端ID不匹配")
	}

	user, err := s.userService.GetUser(ctx, map[string]any{"id": oauthAuthorizationCode.UserID})
	if err != nil {
		return nil, err
	}

	accessTokenString, err := s.jwtManager.GenerateAccessToken(strconv.FormatUint(uint64(user.ID), 10), user.Username, map[string]any{"role": user.Role})
	if err != nil {
		s.logMgr.Error("生成访问令牌失败", "error", err)
		return nil, errors.New("生成访问令牌失败")
	}

	accessToken := &oauthmodels.OAuthAccessToken{
		AccessToken: accessTokenString,
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(s.cfg.OAuth.AccessTokenExpire),
		ClientID:    oauthAuthorizationCode.ClientID,
		Scope:       oauthAuthorizationCode.Scope,
		UserID:      &oauthAuthorizationCode.UserID,
	}

	// 用于在事务中保存 refresh token 字符串
	var refreshTokenString string

	// 使用数据库事务确保授权码、访问令牌和刷新令牌的一致性，当任一步失败时整体回滚
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// 在事务中标记授权码为已使用
		if err := s.oauthAuthorizeService.MarkCodeAsUsedWithTx(ctx, tx, oauthAuthorizationCode.ID); err != nil {
			return err
		}

		// 在事务中保存 access token，获取自增 ID
		if err := s.oauthAccessTokenRepository.CreateWithTx(ctx, tx, accessToken); err != nil {
			s.logMgr.Error("创建OAuth访问令牌失败", "error", err)
			return errors.New("创建OAuth访问令牌失败")
		}

		// 在事务中使用保存后的 accessToken.ID 生成 refresh token
		var genErr error
		refreshTokenString, genErr = s.jwtManager.GenerateRefreshToken(strconv.FormatUint(uint64(user.ID), 10))
		if genErr != nil {
			return genErr
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &oauthdto.ExchangeAccessTokenResponse{
		AccessToken: oauthdto.OAuthAccessTokenResponse{
			AccessToken: accessTokenString,
			ExpiresIn:   int(s.cfg.OAuth.AccessTokenExpire.Seconds()),
		},
		RefreshToken: oauthdto.OAuthRefreshTokenResponse{
			RefreshToken: refreshTokenString,
			ExpiresIn:    int(s.cfg.OAuth.RefreshTokenExpire.Seconds()),
		},
		Scope:     accessToken.Scope,
		TokenType: "Bearer",
	}, nil
}

func (s *OAuthTokenService) RefreshAccessToken(ctx context.Context, form *oauthdto.RefreshAccessTokenForm, clientID, clientSecret string) (*oauthdto.ExchangeAccessTokenResponse, error) {
	// 校验客户端合法性
	oauthClient, err := s.oauthClientService.GetOAuthClient(ctx, map[string]any{"id": clientID, "client_secret": clientSecret})
	if err != nil {
		return nil, err
	}

	// 校验客户端是否支持 refresh_token 授权类型
	if !utils.IsGrantTypeValid("refresh_token", oauthClient.GrantTypes) {
		return nil, errors.New("客户端不支持refresh_token授权类型")
	}

	// 查询刷新令牌
	refreshToken, err := s.oauthRefreshTokenRepository.Get(ctx, map[string]any{"refresh_token": form.RefreshToken})
	if err != nil {
		return nil, err
	}

	// 校验刷新令牌是否已撤销
	if refreshToken.Revoked {
		return nil, errors.New("刷新令牌已撤销")
	}

	// 校验刷新令牌是否已过期
	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("刷新令牌已过期")
	}

	// 校验刷新令牌的客户端ID是否与当前客户端一致
	if refreshToken.ClientID != clientID {
		return nil, errors.New("刷新令牌客户端ID不匹配")
	}

	// 查询用户信息
	user, err := s.userService.GetUser(ctx, map[string]any{"id": refreshToken.UserID})
	if err != nil {
		return nil, err
	}

	// 生成新的访问令牌
	accessTokenString, err := s.jwtManager.RefreshAccessToken(ctx, form.RefreshToken)
	if err != nil {
		s.logMgr.Error("刷新访问令牌失败", "error", err)
		return nil, errors.New("刷新访问令牌失败")
	}

	accessToken := &oauthmodels.OAuthAccessToken{
		AccessToken: accessTokenString,
		TokenType:   "Bearer",
		ExpiresAt:   time.Now().Add(s.cfg.OAuth.AccessTokenExpire),
		ClientID:    refreshToken.ClientID,
		Scope:       refreshToken.Scope,
		UserID:      &refreshToken.UserID,
	}

	// 用于在事务中保存新的 refresh token 字符串
	var newRefreshTokenString string

	// 使用数据库事务确保访问令牌、新刷新令牌和旧刷新令牌撤销的一致性
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// 在事务中保存新的 access token
		if err := s.oauthAccessTokenRepository.CreateWithTx(ctx, tx, accessToken); err != nil {
			s.logMgr.Error("创建OAuth访问令牌失败", "error", err)
			return errors.New("创建OAuth访问令牌失败")
		}

		// 在事务中撤销旧的 refresh token
		if err := s.RevokeRefreshTokenWithTx(ctx, tx, refreshToken.ID); err != nil {
			return err
		}

		// 在事务中生成新的 refresh token
		var genErr error
		newRefreshTokenString, genErr = s.GenerateRefreshTokenWithTx(ctx, tx, accessToken.ID, refreshToken.ClientID, refreshToken.Scope, refreshToken.UserID, user.Username, user.Role)
		if genErr != nil {
			return genErr
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return &oauthdto.ExchangeAccessTokenResponse{
		AccessToken: oauthdto.OAuthAccessTokenResponse{
			AccessToken: accessTokenString,
			ExpiresIn:   int(s.cfg.OAuth.AccessTokenExpire.Seconds()),
		},
		RefreshToken: oauthdto.OAuthRefreshTokenResponse{
			RefreshToken: newRefreshTokenString,
			ExpiresIn:    int(s.cfg.OAuth.RefreshTokenExpire.Seconds()),
		},
		Scope:     accessToken.Scope,
		TokenType: "Bearer",
	}, nil
}

func (s *OAuthTokenService) GenerateRefreshToken(ctx context.Context, accessTokenID uint, clientID string, scope string, userID uint, username string, role string) (string, error) {
	refreshTokenString, err := s.jwtManager.GenerateRefreshToken(strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		s.logMgr.Error("生成刷新令牌失败", "error", err)
		return "", errors.New("生成刷新令牌失败")
	}

	refreshToken := &oauthmodels.OAuthRefreshToken{
		RefreshToken:  refreshTokenString,
		AccessTokenID: accessTokenID,
		ClientID:      clientID,
		Scope:         scope,
		UserID:        userID,
		ExpiresAt:     time.Now().Add(s.cfg.OAuth.RefreshTokenExpire),
	}

	if err := s.oauthRefreshTokenRepository.Create(ctx, refreshToken); err != nil {
		s.logMgr.Error("创建OAuth刷新令牌失败", "error", err)
		return "", errors.New("创建OAuth刷新令牌失败")
	}

	return refreshTokenString, nil
}

// GenerateRefreshTokenWithTx 在事务中生成并保存刷新令牌
func (s *OAuthTokenService) GenerateRefreshTokenWithTx(ctx context.Context, tx *gorm.DB, accessTokenID uint, clientID string, scope string, userID uint, username string, role string) (string, error) {
	refreshTokenString, err := s.jwtManager.GenerateRefreshToken(strconv.FormatUint(uint64(userID), 10))
	if err != nil {
		s.logMgr.Error("生成刷新令牌失败", "error", err)
		return "", errors.New("生成刷新令牌失败")
	}

	refreshToken := &oauthmodels.OAuthRefreshToken{
		RefreshToken:  refreshTokenString,
		AccessTokenID: accessTokenID,
		ClientID:      clientID,
		Scope:         scope,
		UserID:        userID,
		ExpiresAt:     time.Now().Add(s.cfg.OAuth.RefreshTokenExpire),
	}

	if err := s.oauthRefreshTokenRepository.CreateWithTx(ctx, tx, refreshToken); err != nil {
		s.logMgr.Error("创建OAuth刷新令牌失败", "error", err)
		return "", errors.New("创建OAuth刷新令牌失败")
	}

	return refreshTokenString, nil
}

// GetOAuthRefreshToken 根据条件查询刷新令牌
func (s *OAuthTokenService) GetOAuthRefreshToken(ctx context.Context, conds map[string]any) (*oauthmodels.OAuthRefreshToken, error) {
	refreshToken, err := s.oauthRefreshTokenRepository.Get(ctx, conds)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("刷新令牌不存在")
		}
		s.logMgr.Error("获取刷新令牌失败", "error", err, "conds", conds)
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	return refreshToken, nil
}

// RevokeRefreshTokenWithTx 在事务中撤销刷新令牌
func (s *OAuthTokenService) RevokeRefreshTokenWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := tx.WithContext(ctx).Model(&oauthmodels.OAuthRefreshToken{}).Where("id = ?", id).Update("revoked", true).Error; err != nil {
		s.logMgr.Error("撤销刷新令牌失败", "error", err, "id", id)
		return errors.New("撤销刷新令牌失败")
	}
	return nil
}
