package controllers

import (
	"github.com/3086953492/gokit/response"
	"github.com/3086953492/gokit/validator"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctrl *AuthController) RegisterHandler(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, err)
		return
	}

	if err := ctrl.authService.Register(ctx.Request.Context(), &req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建用户成功", nil)
}
