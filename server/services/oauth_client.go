package services

import (
	"context"

	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/logger"
	"go.uber.org/zap"

	"goauth/dto"
	"goauth/models"
	"goauth/repositories"
)

type OAuthClientService struct {
	oauthClientRepository *repositories.OAuthClientRepository
}

func NewOAuthClientService(oauthClientRepository *repositories.OAuthClientRepository) *OAuthClientService {
	return &OAuthClientService{oauthClientRepository: oauthClientRepository}
}

func (s *OAuthClientService) CreateOAuthClient(ctx context.Context, req *dto.CreateOAuthClientRequest) error {
	client := &models.OAuthClient{
		ClientSecret: req.ClientSecret,
		Name:         req.Name,
		Description:  req.Description,
		Logo:         req.Logo,
		RedirectURIs: req.RedirectURIs,
		GrantTypes:   req.GrantTypes,
		Scopes:       req.Scopes,
		Status:       req.Status,
	}
	if err := s.oauthClientRepository.Create(ctx, client); err != nil {
		return errors.Database().Msg("创建OAuth客户端失败").Err(err).Field("client", client).Log()
	}
	logger.Info("创建OAuth客户端成功", zap.Any("client", client))
	return nil
}
