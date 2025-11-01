package services

import (
	"context"
	"fmt"
	"time"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/errors"

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

func (s *UserService) GetUser(ctx context.Context, conds map[string]any) (*models.User, error) {
	user, err := cache.New[models.User]().Key(fmt.Sprintf("%v", conds)).TTL(10*time.Minute).GetOrSet(ctx, func() (*models.User, error) {
		user, err := s.userRepository.Get(ctx, conds)
		if err != nil {
			if errors.IsNotFoundError(err) {
				return nil, err
			}
			return nil, errors.Database().Msg("获取用户失败").Err(err).Field("conds", conds).Log()
		}
		return user, nil
	})
	if err != nil {
		return nil, errors.NotFound().Msg("未找到用户").Err(err).Build()
	}
	return user, err
}
