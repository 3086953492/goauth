package initialize

import (
	"github.com/3086953492/gokit/errors"
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
		return errors.Internal().Msg("注册自定义规则失败").Err(err).Log()
	}

	return nil
}
