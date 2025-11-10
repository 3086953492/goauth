package services

import (
	"context"
	"fmt"
	"time"

	"github.com/3086953492/gokit/cache"
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
	if err := cache.DeleteByPrefix(ctx, "list_oauth_clients:"); err != nil {
		errors.Database().Msg("删除缓存失败").Err(err).Log() // 记录日志，但继续执行
	}
	return nil
}

func (s *OAuthClientService) ListOAuthClients(ctx context.Context, page, pageSize int, conds map[string]any) (*dto.PaginationResponse[dto.OAuthClientListResponse], error) {
	oauthClientsPagination, err := cache.New[dto.PaginationResponse[dto.OAuthClientListResponse]]().Key(fmt.Sprintf("list_oauth_clients:%v", conds)).TTL(10*time.Minute).GetOrSet(ctx, func() (*dto.PaginationResponse[dto.OAuthClientListResponse], error) {
		oauthClients, total, err := s.oauthClientRepository.List(ctx, page, pageSize, conds)
		if err != nil {
			return nil, errors.Database().Msg("获取OAuth客户端列表失败").Err(err).Field("conds", conds).Log()
		}
		oauthClientsResponse := make([]dto.OAuthClientListResponse, len(oauthClients))
		for i, oauthClient := range oauthClients {
			oauthClientsResponse[i] = dto.OAuthClientListResponse{
				ID:       oauthClient.ID,
				Name:     oauthClient.Name,
				Logo:     oauthClient.Logo,
				Status:   oauthClient.Status,
			}
		}
		return &dto.PaginationResponse[dto.OAuthClientListResponse]{
			Items:      oauthClientsResponse,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: int(total / int64(pageSize)),
		}, nil
	})
	if err != nil {
		return nil, errors.NotFound().Msg("未找到OAuth客户端列表").Err(err).Build()
	}
	return oauthClientsPagination, nil
}

func (s *OAuthClientService) GetOAuthClient(ctx context.Context, conds map[string]any) (*dto.OAuthClientDetailResponse, error) {
	oauthClient, err := cache.New[dto.OAuthClientDetailResponse]().Key(fmt.Sprintf("oauth_client:%v", conds)).TTL(10*time.Minute).GetOrSet(ctx, func() (*dto.OAuthClientDetailResponse, error) {
		oauthClient, err := s.oauthClientRepository.Get(ctx, conds)
		if err != nil {
			if errors.IsNotFoundError(err) {
				return nil, err
			}
			return nil, errors.Database().Msg("获取OAuth客户端失败").Err(err).Field("conds", conds).Log()
		}
		return &dto.OAuthClientDetailResponse{
			ID:           oauthClient.ID,
			Name:         oauthClient.Name,
			Description:  oauthClient.Description,
			Logo:         oauthClient.Logo,
			RedirectURIs: oauthClient.RedirectURIs,
			GrantTypes:   oauthClient.GrantTypes,
			Scopes:       oauthClient.Scopes,
			Status:       oauthClient.Status,
			CreatedAt:    oauthClient.CreatedAt,
			UpdatedAt:    oauthClient.UpdatedAt,
		}, nil
	})
	if err != nil {
		return nil, errors.NotFound().Msg("未找到OAuth客户端").Err(err).Build()
	}
	return oauthClient, nil
}