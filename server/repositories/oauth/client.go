package oauthrepositories

import (
	"context"

	"gorm.io/gorm"

	"goauth/models/oauth"
)

// OAuthClientRepository OAuth客户端仓库实现
type OAuthClientRepository struct {
	db *gorm.DB
}

// NewOAuthClientRepository 创建OAuth客户端仓库实例
func NewOAuthClientRepository(db *gorm.DB) *OAuthClientRepository {
	return &OAuthClientRepository{
		db: db,
	}
}

// Create 创建OAuth客户端
func (r *OAuthClientRepository) Create(ctx context.Context, client *oauthmodels.OAuthClient) error {
	return r.db.WithContext(ctx).Create(client).Error
}

// Get 根据传入的条件查询OAuth客户端
func (r *OAuthClientRepository) Get(ctx context.Context, conds map[string]any) (*oauthmodels.OAuthClient, error) {
	var client oauthmodels.OAuthClient
	query := r.db.WithContext(ctx).Model(&oauthmodels.OAuthClient{})

	for key, value := range conds {
		query = query.Where(key, value)
	}

	if err := query.First(&client).Error; err != nil {
		return nil, err
	}

	return &client, nil
}

// Update 更新OAuth客户端信息
func (r *OAuthClientRepository) Update(ctx context.Context, id uint, updates map[string]any) error {
	return r.db.WithContext(ctx).Model(&oauthmodels.OAuthClient{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删除OAuth客户端
func (r *OAuthClientRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&oauthmodels.OAuthClient{}, id).Error
}

// List 分页查询OAuth客户端列表
func (r *OAuthClientRepository) List(ctx context.Context, page, pageSize int, conds map[string]any) ([]oauthmodels.OAuthClient, int64, error) {
	var clients []oauthmodels.OAuthClient
	var total int64

	// 计算总数
	query := r.db.WithContext(ctx).Model(&oauthmodels.OAuthClient{})
	for key, value := range conds {
		query = query.Where(key, value)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Offset(offset).
		Limit(pageSize).
		Find(&clients).Error

	if err != nil {
		return nil, 0, err
	}

	return clients, total, nil
}

