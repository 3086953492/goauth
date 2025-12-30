package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/3086953492/gokit/cache"
	"github.com/3086953492/gokit/logger"
	"github.com/3086953492/gokit/redis"
	"github.com/3086953492/gokit/security/password"
	"github.com/3086953492/gokit/security/subject"
	"github.com/3086953492/gokit/storage"
	"gorm.io/gorm"

	"goauth/dto"
	"goauth/apperrors"
	"goauth/models"
	"goauth/repositories"
	"goauth/utils"
)

// UserService 用户服务实现
type UserService struct {
	userRepository *repositories.UserRepository
	storageManager *storage.Manager
	redisMgr       *redis.Manager
	cacheMgr       *cache.Manager
	logMgr         *logger.Manager
	passwordMgr    *password.Manager
	subjectMgr     *subject.Manager
}

// NewUserService 创建用户服务实例
func NewUserService(userRepository *repositories.UserRepository, storageManager *storage.Manager, redisMgr *redis.Manager, cacheMgr *cache.Manager, logMgr *logger.Manager, passwordMgr *password.Manager, subjectMgr *subject.Manager) *UserService {
	return &UserService{userRepository: userRepository, storageManager: storageManager, redisMgr: redisMgr, cacheMgr: cacheMgr, logMgr: logMgr, passwordMgr: passwordMgr, subjectMgr: subjectMgr}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserForm, avatarFile *utils.FormFileResult) error {

	// 对用户名加锁，防止并发创建相同用户名。
	lockKey := fmt.Sprintf("user:create:%s", req.Username)
	lock := s.redisMgr.NewDistributedLock(lockKey, 10*time.Second)
	if err := lock.Acquire(ctx); err != nil {
		return apperrors.ErrUserSystemBusy
	}
	defer lock.Release(ctx)

	hashedPassword, err := s.passwordMgr.Hash(req.Password)
	if err != nil {
		s.logMgr.Error("密码哈希失败", "error", err)
		return apperrors.ErrUserPasswordHashFailed
	}

	var avatarURL string
	if avatarFile != nil {
		f, err := avatarFile.FileHeader.Open()
		if err != nil {
			s.logMgr.Error("文件读取失败", "error", err)
			return apperrors.ErrUserFileReadFailed
		}
		defer f.Close()

		objectKey := time.Now().Format("2006/01/02") + "/" + avatarFile.Filename
		meta, err := s.storageManager.Upload(ctx, objectKey, f, storage.WithContentType(avatarFile.ContentType))
		if err != nil {
			s.logMgr.Error("头像上传失败", "error", err)
			return apperrors.ErrUserAvatarUploadFailed
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

	// 使用事务创建用户并更新 subject
	err = s.userRepository.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.CreateWithTx(ctx, tx, user); err != nil {
			s.logMgr.Error("创建用户失败", "error", err, "user", user)
			return apperrors.ErrUserCreateFailed
		}

		subject, err := s.subjectMgr.Sub(strconv.FormatUint(uint64(user.ID), 10))
		if err != nil {
			s.logMgr.Error("生成用户标识失败", "error", err)
			return apperrors.ErrUserSubjectGenFailed
		}

		if err := s.userRepository.UpdateWithTx(ctx, tx, user.ID, map[string]any{"subject": subject}); err != nil {
			s.logMgr.Error("更新用户标识失败", "error", err)
			return apperrors.ErrUserSubjectUpdateFailed
		}

		user.Subject = subject
		return nil
	})
	if err != nil {
		return err
	}

	s.logMgr.Info("用户注册成功", "userID", user.ID)

	return nil
}

func (s *UserService) GetUser(ctx context.Context, conds map[string]any) (*models.User, error) {
	user, err := cache.NewBuilder[models.User](s.cacheMgr).KeyWithConds("user", conds).TTL(10*time.Minute).GetOrSet(ctx, func() (*models.User, error) {
		user, err := s.userRepository.Get(ctx, conds)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			s.logMgr.Error("获取用户失败", "error", err, "conds", conds)
			return nil, apperrors.ErrUserSystemBusy
		}
		return user, nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		s.logMgr.Error("获取用户失败", "error", err, "conds", conds)
		return nil, apperrors.ErrUserSystemBusy
	}
	return user, nil
}

// GetUserWithDeleted 获取用户（包含已软删除的记录），用于唯一性验证
func (s *UserService) GetUserWithDeleted(ctx context.Context, conds map[string]any) (*models.User, error) {
	user, err := s.userRepository.GetWithDeleted(ctx, conds)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}
		s.logMgr.Error("获取用户失败", "error", err, "conds", conds)
		return nil, apperrors.ErrUserSystemBusy
	}
	return user, nil
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
		hashedPassword, err := s.passwordMgr.Hash(user.Password)
		if err != nil {
			s.logMgr.Error("密码哈希失败", "error", err)
			return apperrors.ErrUserPasswordHashFailed
		}
		updates["password"] = hashedPassword
	}

	if avatarFile != nil {
		f, err := avatarFile.FileHeader.Open()
		if err != nil {
			s.logMgr.Error("文件读取失败", "error", err)
			return apperrors.ErrUserFileReadFailed
		}
		defer f.Close()

		objectKey := time.Now().Format("2006/01/02") + "/" + avatarFile.Filename
		meta, err := s.storageManager.Upload(ctx, objectKey, f, storage.WithContentType(avatarFile.ContentType))
		if err != nil {
			s.logMgr.Error("头像上传失败", "error", err)
			return apperrors.ErrUserAvatarUploadFailed
		}
		if err := s.storageManager.DeleteByURL(ctx, existingUser.Avatar); err != nil {
			s.logMgr.Warn("旧头像清理失败", "error", err)
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
		s.logMgr.Error("更新用户失败", "error", err, "user", user)
		return apperrors.ErrUserUpdateFailed
	}

	// 删除相关缓存
	if err := s.cacheMgr.DeleteByContainsList(ctx, "user", []map[string]any{{"id": userID}, {"nickname": existingUser.Nickname}, {"username": existingUser.Username}, {"nickname": user.Nickname}}); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err, "user_id", userID)
	}
	if err := s.cacheMgr.DeleteByPrefix(ctx, "list_users:"); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err)
	}

	return nil
}

func (s *UserService) ListUsers(ctx context.Context, page, pageSize int, conds map[string]any) (*dto.PaginationResponse[dto.UserListResponse], error) {
	usersPagination, err := cache.NewBuilder[dto.PaginationResponse[dto.UserListResponse]](s.cacheMgr).Key(fmt.Sprintf("list_users:%v", conds)).TTL(10*time.Minute).GetOrSet(ctx, func() (*dto.PaginationResponse[dto.UserListResponse], error) {
		users, total, err := s.userRepository.List(ctx, page, pageSize, conds)
		if err != nil {
			s.logMgr.Error("获取用户列表失败", "error", err, "conds", conds)
			return nil, apperrors.ErrUserListFailed
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
		return nil, err
	}
	return usersPagination, nil
}

func (s *UserService) DeleteUser(ctx context.Context, userID uint) error {
	lockKey := fmt.Sprintf("user:delete:%v", userID)
	lock := s.redisMgr.NewDistributedLock(lockKey, 10*time.Second)
	if err := lock.Acquire(ctx); err != nil {
		return apperrors.ErrUserSystemBusy
	}
	defer lock.Release(ctx)

	user, err := s.GetUser(ctx, map[string]any{"id": userID})
	if err != nil {
		return err
	}

	if err := s.userRepository.Delete(ctx, userID); err != nil {
		s.logMgr.Error("删除用户失败", "error", err, "user_id", userID)
		return apperrors.ErrUserDeleteFailed
	}

	if err := s.cacheMgr.DeleteByContainsList(ctx, "user", []map[string]any{{"id": userID}, {"nickname": user.Nickname}, {"username": user.Username}}); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err, "user_id", userID)
	}
	if err := s.cacheMgr.DeleteByPrefix(ctx, "list_users:"); err != nil {
		s.logMgr.Warn("删除缓存失败", "error", err)
	}
	s.logMgr.Info("用户删除成功", "userID", userID)
	return nil
}

func (s *UserService) ResolveExtra(ctx context.Context, userIDStr string) (map[string]any, error) {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := s.GetUser(ctx, map[string]any{"id": userID})
	if err != nil {
		return nil, err
	}

	return map[string]any{"role": user.Role}, nil
}
