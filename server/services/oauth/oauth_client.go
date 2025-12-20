package oauthservices

import (
	"context"
	"errors"
	"time"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/logger"
	"gorm.io/gorm"

	"goauth/dto"
	"goauth/models"
	"goauth/repositories/oauth"
)

type OAuthClientService struct {
	oauthClientRepository *oauthrepositories.OAuthClientRepository
	cacheMgr              *cache.Manager
	logMgr                *logger.Manager
}

func NewOAuthClientService(oauthClientRepository *oauthrepositories.OAuthClientRepository, cacheMgr *cache.Manager, logMgr *logger.Manager) *OAuthClientService {
	return &OAuthClientService{oauthClientRepository: oauthClientRepository, cacheMgr: cacheMgr, logMgr: logMgr}
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
		s.logMgr.Error("创建OAuth客户端失败", "error", err, "client", client)
		return errors.New("创建OAuth客户端失败")
	}
	s.logMgr.Info("创建OAuth客户端成功", "client", client)
	if err := s.cacheMgr.DeleteByPrefix(ctx, "list_oauth_clients:"); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err)
	}
	return nil
}

func (s *OAuthClientService) ListOAuthClients(ctx context.Context, page, pageSize int, conds map[string]any) (*dto.PaginationResponse[dto.OAuthClientListResponse], error) {
	oauthClientsPagination, err := cache.NewBuilder[dto.PaginationResponse[dto.OAuthClientListResponse]](s.cacheMgr).KeyWithConds("list_oauth_clients", conds).TTL(10*time.Minute).GetOrSet(ctx, func() (*dto.PaginationResponse[dto.OAuthClientListResponse], error) {
		oauthClients, total, err := s.oauthClientRepository.List(ctx, page, pageSize, conds)
		if err != nil {
			s.logMgr.Error("获取OAuth客户端列表失败", "error", err, "conds", conds)
			return nil, errors.New("获取OAuth客户端列表失败")
		}
		oauthClientsResponse := make([]dto.OAuthClientListResponse, len(oauthClients))
		for i, oauthClient := range oauthClients {
			oauthClientsResponse[i] = dto.OAuthClientListResponse{
				ID:     oauthClient.ID,
				Name:   oauthClient.Name,
				Logo:   oauthClient.Logo,
				Status: oauthClient.Status,
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
		return nil, err
	}
	return oauthClientsPagination, nil
}

func (s *OAuthClientService) GetOAuthClient(ctx context.Context, conds map[string]any) (*dto.OAuthClientDetailResponse, error) {
	oauthClient, err := cache.NewBuilder[dto.OAuthClientDetailResponse](s.cacheMgr).KeyWithConds("oauth_client", conds).TTL(10*time.Minute).GetOrSet(ctx, func() (*dto.OAuthClientDetailResponse, error) {
		oauthClient, err := s.oauthClientRepository.Get(ctx, conds)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			s.logMgr.Error("获取OAuth客户端失败", "error", err, "conds", conds)
			return nil, errors.New("系统繁忙，请稍后再试")
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("OAuth客户端不存在")
		}
		return nil, err
	}
	return oauthClient, nil
}

func (s *OAuthClientService) UpdateOAuthClient(ctx context.Context, id uint, req *dto.UpdateOAuthClientRequest) error {
	updates := make(map[string]any)
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.Logo != nil {
		updates["logo"] = req.Logo
	}
	if req.RedirectURIs != nil {
		updates["redirect_uris"] = req.RedirectURIs
	}
	if req.GrantTypes != nil {
		updates["grant_types"] = req.GrantTypes
	}
	if req.Scopes != nil {
		updates["scopes"] = req.Scopes
	}
	if req.Status != nil {
		updates["status"] = req.Status
	}

	if err := s.oauthClientRepository.Update(ctx, id, updates); err != nil {
		s.logMgr.Error("更新OAuth客户端失败", "error", err, "id", id, "updates", updates)
		return errors.New("更新OAuth客户端失败")
	}

	if err := s.cacheMgr.DeleteByPrefix(ctx, "list_oauth_clients:"); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err)
	}
	if err := s.cacheMgr.DeleteByConds(ctx, "oauth_client", map[string]any{"id": id}); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err)
	}

	return nil
}

func (s *OAuthClientService) DeleteOAuthClient(ctx context.Context, id uint) error {
	if err := s.oauthClientRepository.Delete(ctx, id); err != nil {
		s.logMgr.Error("删除OAuth客户端失败", "error", err, "id", id)
		return errors.New("删除OAuth客户端失败")
	}
	if err := s.cacheMgr.DeleteByPrefix(ctx, "list_oauth_clients:"); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err)
	}
	if err := s.cacheMgr.DeleteByConds(ctx, "oauth_client", map[string]any{"id": id}); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err)
	}
	s.logMgr.Info("删除OAuth客户端成功", "id", id)
	return nil
}
