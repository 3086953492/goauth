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

	"goauth/models/oauth"
	"goauth/repositories/oauth"
)

type OAuthRefreshTokenService struct {
	oauthRefreshTokenRepository *oauthrepositories.OAuthRefreshTokenRepository
	jwtManager                  *jwt.Manager
	cfg                         *config.Config
	logMgr                      *logger.Manager
}

func NewOAuthRefreshTokenService(oauthRefreshTokenRepository *oauthrepositories.OAuthRefreshTokenRepository, jwtManager *jwt.Manager, cfg *config.Config, logMgr *logger.Manager) *OAuthRefreshTokenService {
	return &OAuthRefreshTokenService{oauthRefreshTokenRepository: oauthRefreshTokenRepository, jwtManager: jwtManager, cfg: cfg, logMgr: logMgr}
}

func (s *OAuthRefreshTokenService) GenerateRefreshToken(ctx context.Context, accessTokenID uint, clientID string, scope string, userID uint, username string, role string) (string, error) {
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
func (s *OAuthRefreshTokenService) GenerateRefreshTokenWithTx(ctx context.Context, tx *gorm.DB, accessTokenID uint, clientID string, scope string, userID uint, username string, role string) (string, error) {
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
func (s *OAuthRefreshTokenService) GetOAuthRefreshToken(ctx context.Context, conds map[string]any) (*oauthmodels.OAuthRefreshToken, error) {
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
func (s *OAuthRefreshTokenService) RevokeRefreshTokenWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := tx.WithContext(ctx).Model(&oauthmodels.OAuthRefreshToken{}).Where("id = ?", id).Update("revoked", true).Error; err != nil {
		s.logMgr.Error("撤销刷新令牌失败", "error", err, "id", id)
		return errors.New("撤销刷新令牌失败")
	}
	return nil
}
