package oauthservices

import (
	"context"
	"errors"
	"time"

	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/security/random"
	"gorm.io/gorm"

	"goauth/models/oauth"
	"goauth/repositories/oauth"
)

type OAuthAuthorizeService struct {
	oauthAuthorizationCodeRepository *oauthrepositories.OAuthAuthorizationCodeRepository
	oauthClientService               *OAuthClientService
	logMgr                           *logger.Manager
}

func NewOAuthAuthorizeService(oauthAuthorizationCodeRepository *oauthrepositories.OAuthAuthorizationCodeRepository, oauthClientService *OAuthClientService, logMgr *logger.Manager) *OAuthAuthorizeService {
	return &OAuthAuthorizeService{oauthAuthorizationCodeRepository: oauthAuthorizationCodeRepository, oauthClientService: oauthClientService, logMgr: logMgr}
}

func (s *OAuthAuthorizeService) GenerateAuthorizationCode(ctx context.Context, userID uint, clientID string, redirectURI string, scope string) (string, error) {

	codeString, err := random.URLSafe(32)
	if err != nil {
		s.logMgr.Error("生成授权码失败", "error", err)
		return "", errors.New("生成授权码失败")
	}

	client, err := s.oauthClientService.GetOAuthClient(ctx, map[string]any{"id": clientID})
	if err != nil {
		s.logMgr.Error("获取OAuth客户端失败", "error", err)
		return "", errors.New("系统繁忙，请稍后再试")
	}
	

	code := &oauthmodels.OAuthAuthorizationCode{
		Code:        codeString,
		UserID:      userID,
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Scope:       scope,
		ExpiresAt:   time.Now().Add(time.Duration(client.AuthCodeExpire) * time.Second),
	}
	if err := s.oauthAuthorizationCodeRepository.Create(ctx, code); err != nil {
		s.logMgr.Error("创建OAuth授权码失败", "error", err)
		return "", errors.New("创建OAuth授权码失败")
	}

	return codeString, nil
}

func (s *OAuthAuthorizeService) GetOAuthAuthorizationCode(ctx context.Context, conds map[string]any) (*oauthmodels.OAuthAuthorizationCode, error) {
	oauthAuthorizationCode, err := s.oauthAuthorizationCodeRepository.Get(ctx, conds)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	return oauthAuthorizationCode, nil
}

func (s *OAuthAuthorizeService) MarkCodeAsUsed(ctx context.Context, id uint) error {
	if err := s.oauthAuthorizationCodeRepository.MarkAsUsed(ctx, id); err != nil {
		s.logMgr.Error("标记授权码为已使用失败", "error", err)
		return errors.New("标记授权码为已使用失败")
	}
	return nil
}

// MarkAsUsedWithTx 在事务中标记授权码为已使用
func (s *OAuthAuthorizeService) MarkCodeAsUsedWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := s.oauthAuthorizationCodeRepository.MarkAsUsedWithTx(ctx, tx, id); err != nil {
		s.logMgr.Error("标记授权码为已使用失败", "error", err)
		return errors.New("标记授权码为已使用失败")
	}
	return nil
}
