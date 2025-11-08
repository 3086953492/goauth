package services

import (
	"context"
	"fmt"
	"time"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/redis"
	"go.uber.org/zap"

	"goauth/dto"
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

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) error {

	// 对用户名加锁，防止并发创建相同用户名。
	lockKey := fmt.Sprintf("user:create:%s", req.Username)
	lock := redis.NewDistributedLock(lockKey, 10*time.Second)
	if err := lock.Acquire(); err != nil {
		return errors.Internal().Msg("系统繁忙，请稍后再试").Err(err).Field("username", req.Username).Build()
	}
	defer lock.Release()

	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return errors.Internal().Msg("密码哈希失败").Err(err).Log()
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Role:     "user",
	}
	if err := s.userRepository.Create(ctx, user); err != nil {
		return errors.Database().Err(err).Field("user", user).Log()
	}
	logger.Info("用户注册成功", zap.Uint("userID", user.ID))

	return nil
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

func (s *UserService) UpdateUser(ctx context.Context, userID uint, user *dto.UpdateUserRequest) error {

	// 获取更新前的用户信息
	existingUser, err := s.GetUser(ctx, map[string]any{"id": userID})
	if err != nil {
		return err
	}

	// 构建更新字段 map
	updates := make(map[string]any)

	if user.Nickname != "" {
		updates["nickname"] = user.Nickname
	}

	if user.Password != "" {
		hashedPassword, err := crypto.HashPassword(user.Password)
		if err != nil {
			return errors.Internal().Msg("密码哈希失败").Err(err).Log()
		}
		updates["password"] = hashedPassword
	}

	if user.Avatar != nil {
		updates["avatar"] = *user.Avatar
	}

	if user.Status != nil {
		updates["status"] = *user.Status
	}

	if user.Role != "" {
		updates["role"] = user.Role
	}

	// 执行更新
	if err := s.userRepository.Update(ctx, userID, updates); err != nil {
		return errors.Database().Msg("更新用户失败").Err(err).Field("user", user).Log()
	}

	// 删除相关缓存
	if err := cache.DeleteByContainsList(ctx, []string{fmt.Sprintf("id:%v", userID), fmt.Sprintf("nickname:%v", existingUser.Nickname), fmt.Sprintf("username:%v", existingUser.Username),fmt.Sprintf("nickname:%v", user.Nickname)}); err != nil {
		errors.Database().Msg("删除缓存失败").Err(err).Field("user_id", userID).Log() // 记录日志，但继续执行
	}

	return nil
}
