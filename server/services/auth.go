package services

import (
	"context"
	"strconv"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/jwt"

	"goauth/dto"
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
