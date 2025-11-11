package repositories

import (
	"context"
	"time"

	"gorm.io/gorm"

	"goauth/models"
)

// OAuthAuthorizationCodeRepository OAuth授权码仓库实现
type OAuthAuthorizationCodeRepository struct {
	db *gorm.DB
}

// NewOAuthAuthorizationCodeRepository 创建OAuth授权码仓库实例
func NewOAuthAuthorizationCodeRepository(db *gorm.DB) *OAuthAuthorizationCodeRepository {
	return &OAuthAuthorizationCodeRepository{
		db: db,
	}
}

// Create 创建OAuth授权码
func (r *OAuthAuthorizationCodeRepository) Create(ctx context.Context, code *models.OAuthAuthorizationCode) error {
	return r.db.WithContext(ctx).Create(code).Error
}

// Get 根据传入的条件查询OAuth授权码
func (r *OAuthAuthorizationCodeRepository) Get(ctx context.Context, conds map[string]any) (*models.OAuthAuthorizationCode, error) {
	var code models.OAuthAuthorizationCode
	query := r.db.WithContext(ctx).Model(&models.OAuthAuthorizationCode{})

	for key, value := range conds {
		query = query.Where(key, value)
	}

	if err := query.First(&code).Error; err != nil {
		return nil, err
	}

	return &code, nil
}

// Update 更新OAuth授权码信息
func (r *OAuthAuthorizationCodeRepository) Update(ctx context.Context, id uint, updates map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.OAuthAuthorizationCode{}).Where("id = ?", id).Updates(updates).Error
}

// MarkAsUsed 标记授权码为已使用
func (r *OAuthAuthorizationCodeRepository) MarkAsUsed(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.OAuthAuthorizationCode{}).
		Where("id = ?", id).
		Update("used", true).Error
}

// Delete 软删除OAuth授权码
func (r *OAuthAuthorizationCodeRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.OAuthAuthorizationCode{}, id).Error
}

// DeleteExpired 清理过期的授权码
func (r *OAuthAuthorizationCodeRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.OAuthAuthorizationCode{}).Error
}
