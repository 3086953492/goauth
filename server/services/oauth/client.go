package oauthservices

import (
	"context"
	"errors"
	"time"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/logger"
	"gorm.io/gorm"

	"goauth/dto"
	oauthdto "goauth/dto/oauth"
	oauthmodels "goauth/models/oauth"
	oauthrepositories "goauth/repositories/oauth"
)

// 配置字段默认值（单位：秒）
const (
	DefaultAuthCodeExpire     = 300     // 授权码过期时间：5分钟
	DefaultAccessTokenExpire  = 3600    // 访问令牌过期时间：1小时
	DefaultRefreshTokenExpire = 2592000 // 刷新令牌过期时间：30天
	DefaultSubjectLength      = 16      // 用户标识长度
	DefaultSubjectPrefix      = "usr_"  // 用户标识前缀
)

type OAuthClientService struct {
	oauthClientRepository *oauthrepositories.OAuthClientRepository
	cacheMgr              *cache.Manager
	logMgr                *logger.Manager
}

func NewOAuthClientService(oauthClientRepository *oauthrepositories.OAuthClientRepository, cacheMgr *cache.Manager, logMgr *logger.Manager) *OAuthClientService {
	return &OAuthClientService{oauthClientRepository: oauthClientRepository, cacheMgr: cacheMgr, logMgr: logMgr}
}

func (s *OAuthClientService) CreateOAuthClient(ctx context.Context, req *oauthdto.CreateOAuthClientRequest) error {
	// 配置字段应用默认值
	authCodeExpire := DefaultAuthCodeExpire
	if req.AuthCodeExpire != nil {
		authCodeExpire = *req.AuthCodeExpire
	}

	accessTokenExpire := DefaultAccessTokenExpire
	if req.AccessTokenExpire != nil {
		accessTokenExpire = *req.AccessTokenExpire
	}

	refreshTokenExpire := DefaultRefreshTokenExpire
	if req.RefreshTokenExpire != nil {
		refreshTokenExpire = *req.RefreshTokenExpire
	}

	subjectLength := DefaultSubjectLength
	if req.SubjectLength != nil {
		subjectLength = *req.SubjectLength
	}

	subjectPrefix := DefaultSubjectPrefix
	if req.SubjectPrefix != nil {
		subjectPrefix = *req.SubjectPrefix
	}

	client := &oauthmodels.OAuthClient{
		// 密钥字段（必填）
		ClientSecret:       req.ClientSecret,
		AccessTokenSecret:  req.AccessTokenSecret,
		RefreshTokenSecret: req.RefreshTokenSecret,
		SubjectSecret:      req.SubjectSecret,

		// 基本字段
		Name:         req.Name,
		Description:  req.Description,
		Logo:         req.Logo,
		RedirectURIs: req.RedirectURIs,
		GrantTypes:   req.GrantTypes,
		Scopes:       req.Scopes,
		Status:       req.Status,

		// 配置字段（带默认值）
		AuthCodeExpire:     authCodeExpire,
		AccessTokenExpire:  accessTokenExpire,
		RefreshTokenExpire: refreshTokenExpire,
		SubjectLength:      subjectLength,
		SubjectPrefix:      subjectPrefix,
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

func (s *OAuthClientService) ListOAuthClients(ctx context.Context, page, pageSize int, conds map[string]any) (*dto.PaginationResponse[oauthdto.OAuthClientListResponse], error) {
	oauthClientsPagination, err := cache.NewBuilder[dto.PaginationResponse[oauthdto.OAuthClientListResponse]](s.cacheMgr).KeyWithConds("list_oauth_clients", conds).TTL(10*time.Minute).GetOrSet(ctx, func() (*dto.PaginationResponse[oauthdto.OAuthClientListResponse], error) {
		oauthClients, total, err := s.oauthClientRepository.List(ctx, page, pageSize, conds)
		if err != nil {
			s.logMgr.Error("获取OAuth客户端列表失败", "error", err, "conds", conds)
			return nil, errors.New("获取OAuth客户端列表失败")
		}
		oauthClientsResponse := make([]oauthdto.OAuthClientListResponse, len(oauthClients))
		for i, oauthClient := range oauthClients {
			oauthClientsResponse[i] = oauthdto.OAuthClientListResponse{
				ID:     oauthClient.ID,
				Name:   oauthClient.Name,
				Logo:   oauthClient.Logo,
				Status: oauthClient.Status,
			}
		}
		return &dto.PaginationResponse[oauthdto.OAuthClientListResponse]{
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

func (s *OAuthClientService) GetOAuthClient(ctx context.Context, conds map[string]any) (*oauthdto.OAuthClientDetailResponse, error) {
	oauthClient, err := cache.NewBuilder[oauthdto.OAuthClientDetailResponse](s.cacheMgr).KeyWithConds("oauth_client", conds).TTL(10*time.Minute).GetOrSet(ctx, func() (*oauthdto.OAuthClientDetailResponse, error) {
		oauthClient, err := s.oauthClientRepository.Get(ctx, conds)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			s.logMgr.Error("获取OAuth客户端失败", "error", err, "conds", conds)
			return nil, errors.New("系统繁忙，请稍后再试")
		}
		return &oauthdto.OAuthClientDetailResponse{
			ID:           oauthClient.ID,
			Name:         oauthClient.Name,
			Description:  oauthClient.Description,
			Logo:         oauthClient.Logo,
			RedirectURIs: oauthClient.RedirectURIs,
			GrantTypes:   oauthClient.GrantTypes,
			Scopes:       oauthClient.Scopes,
			Status:       oauthClient.Status,

			// 配置字段（不暴露密钥）
			AuthCodeExpire:     oauthClient.AuthCodeExpire,
			AccessTokenExpire:  oauthClient.AccessTokenExpire,
			RefreshTokenExpire: oauthClient.RefreshTokenExpire,
			SubjectLength:      oauthClient.SubjectLength,
			SubjectPrefix:      oauthClient.SubjectPrefix,

			CreatedAt: oauthClient.CreatedAt,
			UpdatedAt: oauthClient.UpdatedAt,
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

func (s *OAuthClientService) UpdateOAuthClient(ctx context.Context, id uint, req *oauthdto.UpdateOAuthClientRequest) error {
	updates := make(map[string]any)

	// 基本字段
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Logo != nil {
		updates["logo"] = *req.Logo
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
		updates["status"] = *req.Status
	}

	// 密钥字段（轮换：非空时更新）
	if req.ClientSecret != nil && *req.ClientSecret != "" {
		updates["client_secret"] = *req.ClientSecret
	}
	if req.AccessTokenSecret != nil && *req.AccessTokenSecret != "" {
		updates["access_token_secret"] = *req.AccessTokenSecret
	}
	if req.RefreshTokenSecret != nil && *req.RefreshTokenSecret != "" {
		updates["refresh_token_secret"] = *req.RefreshTokenSecret
	}
	if req.SubjectSecret != nil && *req.SubjectSecret != "" {
		updates["subject_secret"] = *req.SubjectSecret
	}

	// 配置字段
	if req.AuthCodeExpire != nil {
		updates["auth_code_expire"] = *req.AuthCodeExpire
	}
	if req.AccessTokenExpire != nil {
		updates["access_token_expire"] = *req.AccessTokenExpire
	}
	if req.RefreshTokenExpire != nil {
		updates["refresh_token_expire"] = *req.RefreshTokenExpire
	}
	if req.SubjectLength != nil {
		updates["subject_length"] = *req.SubjectLength
	}
	if req.SubjectPrefix != nil {
		updates["subject_prefix"] = *req.SubjectPrefix
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
