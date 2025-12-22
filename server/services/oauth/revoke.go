package oauthservices

import (
	"context"
	"errors"
	"goauth/models/oauth"
	"goauth/repositories/oauth"

	"github.com/3086953492/gokit/logger"
	"gorm.io/gorm"
)

// OAuthRevokeService 令牌撤销服务（RFC7009）
type OAuthRevokeService struct {
	db                          *gorm.DB
	oauthAccessTokenRepository  *oauthrepositories.OAuthAccessTokenRepository
	oauthRefreshTokenRepository *oauthrepositories.OAuthRefreshTokenRepository
	logMgr                      *logger.Manager
}

// NewOAuthRevokeService 创建令牌撤销服务实例
func NewOAuthRevokeService(
	db *gorm.DB,
	oauthAccessTokenRepository *oauthrepositories.OAuthAccessTokenRepository,
	oauthRefreshTokenRepository *oauthrepositories.OAuthRefreshTokenRepository,
	logMgr *logger.Manager,
) *OAuthRevokeService {
	return &OAuthRevokeService{
		db:                          db,
		oauthAccessTokenRepository:  oauthAccessTokenRepository,
		oauthRefreshTokenRepository: oauthRefreshTokenRepository,
		logMgr:                      logMgr,
	}
}

// RevokeToken 撤销令牌（支持 access_token 和 refresh_token，refresh 撤销时级联撤销关联的 access token）
// 按 RFC7009 约定：token 不存在或不属于该 client 时，无差别返回 nil（200）不泄露信息
func (s *OAuthRevokeService) RevokeToken(ctx context.Context, token string, tokenTypeHint string, clientID string) error {
	// 根据 hint 决定查询顺序
	switch tokenTypeHint {
	case "access_token":
		// 先查 access，再查 refresh
		if s.tryRevokeAccessToken(ctx, token, clientID) {
			return nil
		}
		s.tryRevokeRefreshToken(ctx, token, clientID)
		return nil
	case "refresh_token":
		// 先查 refresh，再查 access
		if s.tryRevokeRefreshToken(ctx, token, clientID) {
			return nil
		}
		s.tryRevokeAccessToken(ctx, token, clientID)
		return nil
	default:
		// hint 为空或其他：优先 refresh（级联更完整）
		if s.tryRevokeRefreshToken(ctx, token, clientID) {
			return nil
		}
		s.tryRevokeAccessToken(ctx, token, clientID)
		return nil
	}
}

// tryRevokeAccessToken 尝试撤销 access token，成功返回 true
func (s *OAuthRevokeService) tryRevokeAccessToken(ctx context.Context, token string, clientID string) bool {
	accessToken, err := s.oauthAccessTokenRepository.Get(ctx, map[string]any{"access_token": token})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logMgr.Error("查询access token失败", "error", err)
		}
		return false
	}
	// 校验 client 归属
	if accessToken.ClientID != clientID {
		// 不属于该 client，按 RFC7009 无差别处理
		return true
	}
	// 已撤销则直接返回
	if accessToken.Revoked {
		return true
	}
	// 撤销
	if err := s.oauthAccessTokenRepository.Update(ctx, accessToken.ID, map[string]any{"revoked": true}); err != nil {
		s.logMgr.Error("撤销access token失败", "error", err, "id", accessToken.ID)
		return false
	}
	return true
}

// tryRevokeRefreshToken 尝试撤销 refresh token，并级联撤销关联的 access token，成功返回 true
func (s *OAuthRevokeService) tryRevokeRefreshToken(ctx context.Context, token string, clientID string) bool {
	refreshToken, err := s.oauthRefreshTokenRepository.Get(ctx, map[string]any{"refresh_token": token})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logMgr.Error("查询refresh token失败", "error", err)
		}
		return false
	}
	// 校验 client 归属
	if refreshToken.ClientID != clientID {
		return true
	}
	// 已撤销则直接返回
	if refreshToken.Revoked {
		return true
	}
	// 使用事务级联撤销
	txErr := s.db.Transaction(func(tx *gorm.DB) error {
		// 撤销 refresh token
		if err := tx.WithContext(ctx).Model(&oauthmodels.OAuthRefreshToken{}).Where("id = ?", refreshToken.ID).Update("revoked", true).Error; err != nil {
			return err
		}
		// 级联撤销关联的 access token
		if refreshToken.AccessTokenID != 0 {
			if err := tx.WithContext(ctx).Model(&oauthmodels.OAuthAccessToken{}).Where("id = ?", refreshToken.AccessTokenID).Update("revoked", true).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if txErr != nil {
		s.logMgr.Error("撤销refresh token失败", "error", txErr, "id", refreshToken.ID)
		return false
	}
	return true
}
