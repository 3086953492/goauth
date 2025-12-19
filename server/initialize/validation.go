package initialize

import (
	"errors"

	"github.com/3086953492/gokit/validator"
)

func RegisterValidations(container *Container) error {

	if err := container.ValidatorManager.RegisterRules([]validator.Rule{
		{
			Tag:     "username_unique",
			Message: "用户名（{field}） “{value}” 已被占用",
			Func:    container.UserValidator.UsernameUnique,
		},
	}); err != nil {
		return errors.New("注册自定义规则失败")
	}

	return nil
}
