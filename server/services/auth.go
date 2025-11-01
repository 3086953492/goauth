package services

import (
	"context"
	"fmt"
	"time"

	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/redis"
	"go.uber.org/zap"

	"goauth/dto"
	"goauth/models"
	"goauth/repositories"
)

// AuthService 授权服务实现
type AuthService struct {
	userRepository *repositories.UserRepository
	userService    *UserService
}

// NewAuthService 创建授权服务实例
func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository, userService: NewUserService(userRepository)}
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) error {

	// 对用户名加锁，防止并发创建相同用户名。
	lockKey := fmt.Sprintf("user:register:%s", req.Username)
	lock := redis.NewDistributedLock(lockKey, 10*time.Second)
	if err := lock.Acquire(); err != nil {
		return errors.Internal().Msg("系统繁忙，请稍后再试").Err(err).Field("username", req.Username).Build()
	}
	defer lock.Release()

	user, err := s.userService.GetUser(ctx, map[string]any{"username": req.Username})
	if err != nil && !errors.IsNotFoundError(err) {
		return err
	}
	if user != nil {
		return errors.Validation().Msg("当前用户已被注册").Field("username", req.Username).Build()
	}

	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return errors.Internal().Msg("密码哈希失败").Err(err).Log()
	}

	user = &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   0, //默认禁用
		Role:     "user",
	}
	if err := s.userRepository.Create(ctx, user); err != nil {
		return errors.Database().Err(err).Field("user", user).Log()
	}
	logger.Info("用户注册成功", zap.Uint("userID", user.ID))

	return nil
}
