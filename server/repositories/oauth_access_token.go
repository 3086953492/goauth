package repositories

import (
	"context"

	"gorm.io/gorm"

	"goauth/models"
)

// OAuthAccessTokenRepository OAuth访问令牌仓库实现
type OAuthAccessTokenRepository struct {
	db *gorm.DB
}

// NewOAuthAccessTokenRepository 创建OAuth访问令牌仓库实例
func NewOAuthAccessTokenRepository(db *gorm.DB) *OAuthAccessTokenRepository {
	return &OAuthAccessTokenRepository{
		db: db,
	}
}

// Create 创建OAuth访问令牌
func (r *OAuthAccessTokenRepository) Create(ctx context.Context, token *models.OAuthAccessToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

// CreateWithTx 在事务中创建OAuth访问令牌
func (r *OAuthAccessTokenRepository) CreateWithTx(ctx context.Context, tx *gorm.DB, token *models.OAuthAccessToken) error {
	return tx.WithContext(ctx).Create(token).Error
}

// Get 根据传入的条件查询OAuth访问令牌
func (r *OAuthAccessTokenRepository) Get(ctx context.Context, conds map[string]any) (*models.OAuthAccessToken, error) {
	var token models.OAuthAccessToken
	query := r.db.WithContext(ctx).Model(&models.OAuthAccessToken{})

	for key, value := range conds {
		query = query.Where(key, value)
	}

	if err := query.First(&token).Error; err != nil {
		return nil, err
	}

	return &token, nil
}

// Update 更新OAuth访问令牌信息
func (r *OAuthAccessTokenRepository) Update(ctx context.Context, id uint, updates map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.OAuthAccessToken{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删除OAuth访问令牌
func (r *OAuthAccessTokenRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.OAuthAccessToken{}, id).Error
}

// List 分页查询OAuth访问令牌列表
func (r *OAuthAccessTokenRepository) List(ctx context.Context, page, pageSize int, conds map[string]any) ([]models.OAuthAccessToken, int64, error) {
	var tokens []models.OAuthAccessToken
	var total int64

	// 计算总数
	query := r.db.WithContext(ctx).Model(&models.OAuthAccessToken{})
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
