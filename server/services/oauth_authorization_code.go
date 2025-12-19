package services

import (
	"context"
	"errors"
	"time"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/logger"
	"gorm.io/gorm"

	"goauth/models"
	"goauth/repositories"
)

type OAuthAuthorizationCodeService struct {
	oauthAuthorizationCodeRepository *repositories.OAuthAuthorizationCodeRepository
	oauthClientService               *OAuthClientService
	cfg                              *config.Config
	logMgr                           *logger.Manager
}

func NewOAuthAuthorizationCodeService(oauthAuthorizationCodeRepository *repositories.OAuthAuthorizationCodeRepository, oauthClientService *OAuthClientService, cfg *config.Config, logMgr *logger.Manager) *OAuthAuthorizationCodeService {
	return &OAuthAuthorizationCodeService{oauthAuthorizationCodeRepository: oauthAuthorizationCodeRepository, oauthClientService: oauthClientService, cfg: cfg, logMgr: logMgr}
}

func (s *OAuthAuthorizationCodeService) GenerateAuthorizationCode(ctx context.Context, userID uint, clientID string, redirectURI string, scope string) (string, error) {

	codeString, err := crypto.GenerateAuthorizationCode(32)
	if err != nil {
		s.logMgr.Error("生成授权码失败", "error", err)
		return "", errors.New("生成授权码失败")
	}

	code := &models.OAuthAuthorizationCode{
		Code:        codeString,
		UserID:      userID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(s.cfg.OAuth.AuthCodeExpire),
	}
	if err := s.oauthAuthorizationCodeRepository.Create(ctx, code); err != nil {
		s.logMgr.Error("创建OAuth授权码失败", "error", err)
		return "", errors.New("创建OAuth授权码失败")
	}

	return codeString, nil
}

func (s *OAuthAuthorizationCodeService) GetOAuthAuthorizationCode(ctx context.Context, conds map[string]any) (*models.OAuthAuthorizationCode, error) {
	oauthAuthorizationCode, err := s.oauthAuthorizationCodeRepository.Get(ctx, conds)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	return oauthAuthorizationCode, nil
}

func (s *OAuthAuthorizationCodeService) MarkAsUsed(ctx context.Context, id uint) error {
	if err := s.oauthAuthorizationCodeRepository.MarkAsUsed(ctx, id); err != nil {
		s.logMgr.Error("标记授权码为已使用失败", "error", err)
		return errors.New("标记授权码为已使用失败")
	}
	return nil
}

// MarkAsUsedWithTx 在事务中标记授权码为已使用
func (s *OAuthAuthorizationCodeService) MarkAsUsedWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := s.oauthAuthorizationCodeRepository.MarkAsUsedWithTx(ctx, tx, id); err != nil {
		s.logMgr.Error("标记授权码为已使用失败", "error", err)
		return errors.New("标记授权码为已使用失败")
	}
	return nil
}
