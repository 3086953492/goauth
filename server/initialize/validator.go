package initialize

import (
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/validator"

	"goauth/validations"
)

func InitValidator(container *Container) {
	v := validator.New()
	userValidators := validations.NewUserValidators(container.UserService)

	// 注册所有自定义规则
	if err := v.RegisterRules([]validator.Rule{
		{
			Tag:     "username_unique",
			Message: "用户名 {value} 已被占用",
			Func:    userValidators.UsernameUnique,
		},
	}); err != nil {
		errors.Internal().Msg("注册自定义规则失败").Err(err).Log()
		return
	}

	validator.SetDefaultValidator(v)
}
