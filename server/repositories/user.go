package repositories

import (
	"context"

	"gorm.io/gorm"

	"goauth/models"
)

// UserRepository 用户仓库实现
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓库实例
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Get 根据传入的条件查询用户
func (r *UserRepository) Get(ctx context.Context, conds map[string]any) (*models.User, error) {
	var user models.User
	query := r.db.WithContext(ctx).Model(&models.User{})

	for key, value := range conds {
		query = query.Where(key, value)
	}

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(ctx context.Context, id uint, updates map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 软删除用户
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List 分页查询用户列表
func (r *UserRepository) List(ctx context.Context, page, pageSize int, conds map[string]any) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 计算总数
	query := r.db.WithContext(ctx).Model(&models.User{})
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
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
