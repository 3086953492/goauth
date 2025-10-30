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

// GetByID 根据ID查询用户
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名查询用户
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 软删除用户
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List 分页查询用户列表
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 计算总数
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ExistsByUsername 检查用户名是否已存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetByStatus 根据状态查询用户
func (r *UserRepository) GetByStatus(ctx context.Context, status int, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 计算总数
	query := r.db.WithContext(ctx).Model(&models.User{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

