package controllers

import (
	"github.com/3086953492/gokit/response"
	"github.com/3086953492/gokit/validator"
	"github.com/gin-gonic/gin"

	"goauth/models"
	"goauth/services"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) CreateUserHandler(ctx *gin.Context) {
	var req models.CreateUserRequest
	if err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := ctrl.userService.CreateUser(ctx.Request.Context(), &req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建用户成功", nil)
}
