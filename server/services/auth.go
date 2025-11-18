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

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, *dto.AuthRefreshTokenResponse, error) {

	user, err := s.userRepository.Get(ctx, map[string]any{"username": req.Username})
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, nil, errors.InvalidInput().Msg("账号或密码错误").Build()
		}
		return nil, nil, errors.Database().Msg("系统繁忙，请稍后再试").Err(err).Field("username", req.Username).Log()
	}

	if user.Status == 0 {
		return nil, nil, errors.Forbidden().Msg("账号未激活或已禁用，请联系管理员").Build()
	}

	if !crypto.VerifyPassword(user.Password, req.Password) {
		return nil, nil, errors.InvalidInput().Msg("账号或密码错误").Build()
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	accessToken, err := jwt.GenerateToken(userID, user.Username, map[string]any{"role": user.Role})
	if err != nil {
		return nil, nil, errors.Internal().Msg("生成访问令牌失败").Err(err).Log()
	}

	refreshToken, err := jwt.GenerateRefreshToken(userID)
	if err != nil {
		return nil, nil, errors.Internal().Msg("生成刷新令牌失败").Err(err).Log()
	}

	return &dto.LoginResponse{
			User: &dto.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Nickname:  user.Nickname,
				Avatar:    user.Avatar,
				Status:    user.Status,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
			AccessToken: &dto.AuthAccessTokenResponse{
				AccessToken: accessToken,
				ExpiresIn:   int(config.GetGlobalConfig().JWT.Expire.Seconds()),
			},
		}, &dto.AuthRefreshTokenResponse{
			RefreshToken: refreshToken,
			ExpiresIn:    int(config.GetGlobalConfig().JWT.RefreshExpire.Seconds()),
		}, nil
}


func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthAccessTokenResponse, error) {

	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, errors.Unauthorized().Msg("刷新令牌验证失败").Err(err).Build()
	}

	if claims.TokenType != jwt.RefreshToken {
		return nil, errors.Unauthorized().Msg("无效的刷新令牌").Build()
	}

	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		return nil, errors.Unauthorized().Msg("用户ID格式错误").Build()
	}

	_, err = s.userRepository.Get(ctx, map[string]any{"id": userID})
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.NotFound().Msg("用户不存在").Build()
		}
		return nil, errors.Database().Msg("系统繁忙，请稍后再试").Err(err).Field("user_id", userID).Log()
	}

	accessToken, err := jwt.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, errors.InvalidInput().Msg("刷新令牌验证失败").Err(err).Log()
	}

	return &dto.AuthAccessTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(config.GetGlobalConfig().JWT.Expire.Seconds()),
	}, nil
}