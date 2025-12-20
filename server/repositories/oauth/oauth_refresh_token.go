package oauthrepositories

import (
	"context"

	"gorm.io/gorm"

	"goauth/models"
)

// OAuthRefreshTokenRepository OAuth刷新令牌仓库实现
type OAuthRefreshTokenRepository struct {
	db *gorm.DB
}

// NewOAuthRefreshTokenRepository 创建OAuth刷新令牌仓库实例
func NewOAuthRefreshTokenRepository(db *gorm.DB) *OAuthRefreshTokenRepository {
	return &OAuthRefreshTokenRepository{
		db: db,
	}
}

// Create 创建OAuth刷新令牌
func (r *OAuthRefreshTokenRepository) Create(ctx context.Context, token *models.OAuthRefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// CreateWithTx 在事务中创建OAuth刷新令牌
func (r *OAuthRefreshTokenRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, token *models.OAuthRefreshToken) error {
	return tx.WithContext(ctx).Create(token).Error
}

// Get 根据传入的条件查询OAuth刷新令牌
func (r *OAuthRefreshTokenRepository) Get(ctx context.Context, conds map[string]any) (*models.OAuthRefreshToken, error) {
	var token models.OAuthRefreshToken
	query := r.db.WithContext(ctx).Model(&models.OAuthRefreshToken{})

	for key, value := range conds {
		query = query.Where(key, value)
	}

	if err := query.First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

// Update 更新OAuth刷新令牌信息
func (r *OAuthRefreshTokenRepository) Update(ctx context.Context, id uint, updates map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.OAuthRefreshToken{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删除OAuth刷新令牌
func (r *OAuthRefreshTokenRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.OAuthRefreshToken{}, id).Error
}

// List 分页查询OAuth刷新令牌列表
func (r *OAuthRefreshTokenRepository) List(ctx context.Context, page, pageSize int, conds map[string]any) ([]models.OAuthRefreshToken, int64, error) {
	var tokens []models.OAuthRefreshToken
	var total int64

	// 计算总数
	query := r.db.WithContext(ctx).Model(&models.OAuthRefreshToken{})
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
		Find(&tokens).Error

	if err != nil {
		return nil, 0, err
	}

	return tokens, total, nil
}

