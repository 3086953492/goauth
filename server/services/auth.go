package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/jwt"
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

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := s.userRepository.Get(ctx, map[string]any{"username": req.Username})
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.InvalidInput().Msg("账号或密码错误").Build()
		}
		return nil, errors.Database().Msg("系统繁忙，请稍后再试").Err(err).Field("username", req.Username).Log()
	}

	if user.Status == 0 {
		return nil, errors.Forbidden().Msg("账号未激活或已禁用，请联系管理员").Build()
	}

	if !crypto.VerifyPassword(user.Password, req.Password) {
		return nil, errors.InvalidInput().Msg("账号或密码错误").Build()
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	accessToken, err := jwt.GenerateToken(userID, user.Username, map[string]any{"role": user.Role})
	if err != nil {
		return nil, errors.Internal().Msg("生成访问令牌失败").Err(err).Log()
	}

	refreshToken, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, errors.Internal().Msg("生成刷新令牌失败").Err(err).Log()
	}

	return &dto.LoginResponse{
		User: &dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Status:   user.Status,
			Role:     user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Token: &dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    int(config.GetGlobalConfig().JWT.Expire.Seconds()),
		},
	}, nil
}
