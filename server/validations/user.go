package validations

import (
	"context"

	"github.com/3086953492/gokit/validator"

	"goauth/services"
)

// UserValidators 用户验证器实现
type UserValidators struct {
	userService *services.UserService
}

// NewUserValidators 创建用户验证器实例
func NewUserValidators(userService *services.UserService) *UserValidators {
	return &UserValidators{userService: userService}
}

// UsernameUnique 用户名唯一性验证 -> username_unique
// 使用 GetUserWithDeleted 检查包括已删除用户在内的所有记录，确保唯一性约束与数据库一致
func (v *UserValidators) UsernameUnique(fl validator.FieldLevel) bool {
	username := fl.Field().Interface().(string)
	if username == "" {
		return true
	}
	_, err := v.userService.GetUserWithDeleted(context.Background(), map[string]any{"username": username})
	return err != nil
}
