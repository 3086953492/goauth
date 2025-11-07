package controllers

import (
	"github.com/3086953492/gokit/errors"
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
	if result, err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}

	if err := ctrl.authService.Register(ctx.Request.Context(), &req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建用户成功", nil)
}

func (ctrl *AuthController) LoginHandler(ctx *gin.Context) {
	var req dto.LoginRequest
	if result, err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}
	loginResponse, err := ctrl.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "登录成功", loginResponse)
}