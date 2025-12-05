package services

import (
	"context"
	"time"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/errors"
	"gorm.io/gorm"

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

func (s *OAuthAuthorizationCodeService) GenerateAuthorizationCode(ctx context.Context, userID uint, clientID string, redirectURI string, scope string) (string, error) {

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
		ExpiresAt:   time.Now().Add(config.GetGlobalConfig().OAuth.AuthCodeExpire),
	}
	if err := s.oauthAuthorizationCodeRepository.Create(ctx, code); err != nil {
		return "", errors.Database().Msg("创建OAuth授权码失败").Err(err).Field("code", code).Log()
	}

	return codeString, nil
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

// MarkAsUsedWithTx 在事务中标记授权码为已使用
func (s *OAuthAuthorizationCodeService) MarkAsUsedWithTx(ctx context.Context, tx *gorm.DB, id uint) error {
	if err := s.oauthAuthorizationCodeRepository.MarkAsUsedWithTx(ctx, tx, id); err != nil {
		return errors.Database().Msg("标记授权码为已使用失败").Err(err).Log()
	}
	return nil
}