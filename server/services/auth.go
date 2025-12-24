package services

import (
	"context"
	"errors"
	"strconv"

	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/jwt"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/security/password"
	"gorm.io/gorm"

	"goauth/dto"
	"goauth/repositories"
)

// AuthService 授权服务实现
type AuthService struct {
	userRepository *repositories.UserRepository
	userService    *UserService
	logMgr         *logger.Manager
	jwtManager     *jwt.Manager
	passwordMgr    *password.Manager
	cfg            *config.Config
}

// NewAuthService 创建授权服务实例
func NewAuthService(userRepository *repositories.UserRepository, userService *UserService, logMgr *logger.Manager, jwtManager *jwt.Manager, passwordMgr *password.Manager, cfg *config.Config) *AuthService {
	return &AuthService{userRepository: userRepository, userService: userService, logMgr: logMgr, jwtManager: jwtManager, passwordMgr: passwordMgr, cfg: cfg}
}

func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (accessToken string, accessTokenExpire int, refreshToken string, refreshTokenExpire int, userResp *dto.UserResponse, err error) {

	user, err := s.userRepository.Get(ctx, map[string]any{"username": req.Username})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logMgr.Error("获取用户失败", "error", err)
			return "", 0, "", 0, nil, errors.New("系统繁忙，请稍后再试")
		}
		return "", 0, "", 0, nil, errors.New("账号或密码错误")
	}

	if user.Status == 0 {
		return "", 0, "", 0, nil, errors.New("账号未激活或已禁用，请联系管理员")
	}

	if err := s.passwordMgr.Compare(user.Password, req.Password); err != nil {
		if errors.Is(err, password.ErrMismatch) {
			return "", 0, "", 0, nil, errors.New("账号或密码错误")
		}
		s.logMgr.Error("密码验证失败", "error", err)
		return "", 0, "", 0, nil, errors.New("系统繁忙，请稍后再试")
	}

	userID := strconv.FormatUint(uint64(user.ID), 10)
	accessToken, err = s.jwtManager.GenerateAccessToken(userID, map[string]any{"role": user.Role})
	if err != nil {
		s.logMgr.Error("生成访问令牌失败", "error", err)
		return "", 0, "", 0, nil, errors.New("生成访问令牌失败")
	}

	refreshToken, err = s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		s.logMgr.Error("生成刷新令牌失败", "error", err)
		return "", 0, "", 0, nil, errors.New("生成刷新令牌失败")
	}

	return accessToken, int(s.cfg.AuthToken.AccessTokenExpire.Seconds()), refreshToken, int(s.cfg.AuthToken.RefreshTokenExpire.Seconds()), &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (accessToken string, accessTokenExpire int, err error) {

	claims, err := s.jwtManager.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", 0, errors.New("无效的刷新令牌")
	}

	if claims.TokenType != jwt.RefreshToken {
		return "", 0, errors.New("无效的刷新令牌")
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return "", 0, errors.New("用户ID格式错误")
	}

	_, err = s.userRepository.Get(ctx, map[string]any{"id": userID})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logMgr.Error("获取用户失败", "error", err)
			return "", 0, errors.New("系统繁忙，请稍后再试")
		}
		return "", 0, errors.New("用户不存在")
	}

	accessToken, err = s.jwtManager.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		s.logMgr.Error("刷新令牌失败", "error", err)
		return "", 0, errors.New("系统繁忙，请稍后再试")
	}

	return accessToken, int(s.cfg.AuthToken.AccessTokenExpire.Seconds()), nil
}
