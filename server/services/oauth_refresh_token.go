package services

import (
	"context"
	"strconv"
	"time"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/jwt"
	"gorm.io/gorm"

	"goauth/models"
	"goauth/repositories"
)

type OAuthRefreshTokenService struct {
	oauthRefreshTokenRepository *repositories.OAuthRefreshTokenRepository
}

func NewOAuthRefreshTokenService(oauthRefreshTokenRepository *repositories.OAuthRefreshTokenRepository) *OAuthRefreshTokenService {
	return &OAuthRefreshTokenService{oauthRefreshTokenRepository: oauthRefreshTokenRepository}
}

func (s *OAuthRefreshTokenService) GenerateRefreshToken(ctx context.Context, accessTokenID uint, clientID string, scope string, userID uint, username string, role string) (string, error) {
	refreshTokenString, err := jwt.GenerateToken(strconv.FormatUint(uint64(userID), 10), username, map[string]any{"role": role})
	if err != nil {
		return "", errors.Internal().Msg("生成刷新令牌失败").Err(err).Log()
	}

	refreshToken := &models.OAuthRefreshToken{
		RefreshToken:  refreshTokenString,
		AccessTokenID: accessTokenID,
		ClientID:      clientID,
		Scope:         scope,
		UserID:        userID,
		ExpiresAt:     time.Now().Add(config.GetGlobalConfig().OAuth.RefreshTokenExpire),
	}

	if err := s.oauthRefreshTokenRepository.Create(ctx, refreshToken); err != nil {
		return "", errors.Database().Msg("创建OAuth刷新令牌失败").Err(err).Log()
	}

	return refreshTokenString, nil
}

// GenerateRefreshTokenWithTx 在事务中生成并保存刷新令牌
func (s *OAuthRefreshTokenService) GenerateRefreshTokenWithTx(ctx context.Context, tx *gorm.DB, accessTokenID uint, clientID string, scope string, userID uint, username string, role string) (string, error) {
	refreshTokenString, err := jwt.GenerateToken(strconv.FormatUint(uint64(userID), 10), username, map[string]any{"role": role})
	if err != nil {
		return "", errors.Internal().Msg("生成刷新令牌失败").Err(err).Log()
	}

	refreshToken := &models.OAuthRefreshToken{
		RefreshToken:  refreshTokenString,
		AccessTokenID: accessTokenID,
		ClientID:      clientID,
		Scope:         scope,
		UserID:        userID,
		ExpiresAt:     time.Now().Add(config.GetGlobalConfig().OAuth.RefreshTokenExpire),
	}

	if err := s.oauthRefreshTokenRepository.CreateWithTx(ctx, tx, refreshToken); err != nil {
		return "", errors.Database().Msg("创建OAuth刷新令牌失败").Err(err).Log()
	}

	return refreshTokenString, nil
}
