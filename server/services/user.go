package services

import (
	"context"

	"goauth/models"
	"goauth/repositories"
)

// UserService 用户服务实现
type UserService struct {
	userRepository *repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) error {
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Role:     "user",
	}
	return s.userRepository.Create(ctx, user)
}
