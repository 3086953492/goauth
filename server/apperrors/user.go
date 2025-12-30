package apperrors

import "errors"

// 用户服务业务错误定义

var (
	// 通用错误
	ErrUserSystemBusy = errors.New("系统繁忙，请稍后再试")

	// 用户不存在
	ErrUserNotFound = errors.New("用户不存在")

	// 密码相关
	ErrUserPasswordHashFailed = errors.New("密码哈希失败")

	// 文件操作相关
	ErrUserFileReadFailed     = errors.New("文件读取失败")
	ErrUserAvatarUploadFailed = errors.New("头像上传失败")

	// 用户创建相关
	ErrUserCreateFailed        = errors.New("创建用户失败")
	ErrUserSubjectGenFailed    = errors.New("生成用户标识失败")
	ErrUserSubjectUpdateFailed = errors.New("更新用户标识失败")

	// 用户更新相关
	ErrUserUpdateFailed = errors.New("更新用户失败")

	// 用户删除相关
	ErrUserDeleteFailed = errors.New("删除用户失败")

	// 用户列表相关
	ErrUserListFailed = errors.New("获取用户列表失败")
)
