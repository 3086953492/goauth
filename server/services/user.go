package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/crypto"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/redis"
	"github.com/3086953492/gokit/storage"

	"goauth/dto"
	"goauth/models"
	"goauth/repositories"
	"goauth/utils"
)

// UserService 用户服务实现
type UserService struct {
	userRepository *repositories.UserRepository
	storageManager *storage.Manager
	redisMgr *redis.Manager
	cacheMgr *cache.Manager
	logMgr *logger.Manager
}

// NewUserService 创建用户服务实例
func NewUserService(userRepository *repositories.UserRepository, storageManager *storage.Manager, redisMgr *redis.Manager, cacheMgr *cache.Manager, logMgr *logger.Manager) *UserService {
	return &UserService{userRepository: userRepository, storageManager: storageManager, redisMgr: redisMgr, cacheMgr: cacheMgr, logMgr: logMgr}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserForm, avatarFile *utils.FormFileResult) error {

	// 对用户名加锁，防止并发创建相同用户名。
	lockKey := fmt.Sprintf("user:create:%s", req.Username)
	lock := s.redisMgr.NewDistributedLock(lockKey, 10*time.Second)
	if err := lock.Acquire(ctx); err != nil {
		return errors.Internal().Msg("系统繁忙，请稍后再试").Err(err).Field("username", req.Username).Build()
	}
	defer lock.Release(ctx)

	hashedPassword, err := crypto.HashPassword(req.Password)
	if err != nil {
		return errors.Internal().Msg("密码哈希失败").Err(err).Log()
	}

	var avatarURL string
	if avatarFile != nil {
		f, err := avatarFile.FileHeader.Open()
		if err != nil {
			return errors.Internal().Msg("文件读取失败").Err(err).Log()
		}
		defer f.Close()

		objectKey := time.Now().Format("2006/01/02") + "/" + avatarFile.Filename
		meta, err := s.storageManager.Upload(ctx, objectKey, f, storage.WithContentType(avatarFile.ContentType))
		if err != nil {
			return errors.Internal().Msg("头像上传失败").Err(err).Log()
		}
		avatarURL = meta.URL
	}

	user := &models.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Avatar:   avatarURL,
		Role:     "user",
	}
	if err := s.userRepository.Create(ctx, user); err != nil {
		return errors.Database().Err(err).Field("user", user).Log()
	}
	s.logMgr.Info("用户注册成功", "userID", user.ID)

	return nil
}

func (s *UserService) GetUser(ctx context.Context, conds map[string]any) (*models.User, error) {
	user, err := cache.NewBuilder[models.User](s.cacheMgr).KeyWithConds("user", conds).TTL(10*time.Minute).GetOrSet(ctx, func() (*models.User, error) {
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

func (s *UserService) UpdateUser(ctx context.Context, userID uint, user *dto.UpdateUserForm, avatarFile *utils.FormFileResult) error {

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

	if avatarFile != nil {
		f, err := avatarFile.FileHeader.Open()
		if err != nil {
			return errors.Internal().Msg("文件读取失败").Err(err).Log()
		}
		defer f.Close()

		objectKey := time.Now().Format("2006/01/02") + "/" + avatarFile.Filename
		meta, err := s.storageManager.Upload(ctx, objectKey, f, storage.WithContentType(avatarFile.ContentType))
		if err != nil {
			return errors.Internal().Msg("头像上传失败").Err(err).Log()
		}
		if err := s.storageManager.DeleteByURL(ctx, existingUser.Avatar); err != nil {
			errors.Internal().Msg("旧头像清理失败").Err(err).Log()
		}
		updates["avatar"] = meta.URL
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
	if err := s.cacheMgr.DeleteByContainsList(ctx, "user", []map[string]any{{"id": userID}, {"nickname": existingUser.Nickname}, {"username": existingUser.Username}, {"nickname": user.Nickname}}); err != nil {
		errors.Database().Msg("删除缓存失败").Err(err).Field("user_id", userID).Log() // 记录日志，但继续执行
	}
	if err := s.cacheMgr.DeleteByPrefix(ctx, "list_users:"); err != nil {
		errors.Database().Msg("删除缓存失败").Err(err).Log() // 记录日志，但继续执行
	}

	return nil
}

func (s *UserService) ListUsers(ctx context.Context, page, pageSize int, conds map[string]any) (*dto.PaginationResponse[dto.UserListResponse], error) {
	usersPagination, err := cache.NewBuilder[dto.PaginationResponse[dto.UserListResponse]](s.cacheMgr).Key(fmt.Sprintf("list_users:%v", conds)).TTL(10*time.Minute).GetOrSet(ctx, func() (*dto.PaginationResponse[dto.UserListResponse], error) {
		users, total, err := s.userRepository.List(ctx, page, pageSize, conds)
		if err != nil {
			return nil, errors.Database().Msg("获取用户列表失败").Err(err).Field("conds", conds).Log()
		}
		usersResponse := make([]dto.UserListResponse, len(users))
		for i, user := range users {
			usersResponse[i] = dto.UserListResponse{
				ID:       user.ID,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
				Status:   user.Status,
				Role:     user.Role,
			}
		}
		return &dto.PaginationResponse[dto.UserListResponse]{
			Items:      usersResponse,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: int(total / int64(pageSize)),
		}, nil
	})
	if err != nil {
		return nil, errors.NotFound().Msg("未找到用户列表").Err(err).Build()
	}
	return usersPagination, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	lockKey := fmt.Sprintf("user:delete:%v", userID)
	lock := s.redisMgr.NewDistributedLock(lockKey, 10*time.Second)
	if err := lock.Acquire(ctx); err != nil {
		return errors.Internal().Msg("系统繁忙，请稍后再试").Err(err).Field("user_id", userID).Build()
	}
	defer lock.Release(ctx)

	user, err := s.GetUser(ctx, map[string]any{"id": userID})
	if err != nil {
		return err
	}

	if err := s.userRepository.Delete(ctx, userID); err != nil {
		return errors.Database().Msg("删除用户失败").Err(err).Field("user_id", userID).Log()
	}

	if err := s.cacheMgr.DeleteByContainsList(ctx, "user", []map[string]any{{"id": userID}, {"nickname": user.Nickname}, {"username": user.Username}}); err != nil {
		errors.Database().Msg("删除缓存失败").Err(err).Field("user_id", userID).Log() // 记录日志，但继续执行
	}
	if err := s.cacheMgr.DeleteByPrefix(ctx, "list_users:"); err != nil {
		errors.Database().Msg("删除缓存失败").Err(err).Log() // 记录日志，但继续执行
	}
	s.logMgr.Info("用户删除成功", "userID", userID)
	return nil
}

func (s *UserService) ResolveExtra(ctx context.Context, userIDStr string) (string, map[string]any, error) {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return "", nil, err
	}

	user, err := s.GetUser(ctx, map[string]any{"id": userID})
	if err != nil {
		return "", nil, err
	}

	return user.Username, map[string]any{"role": user.Role}, nil
}
