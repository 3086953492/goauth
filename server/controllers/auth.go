package controllers

import (
	"github.com/3086953492/gokit/config"
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

func (ctrl *AuthController) LoginHandler(ctx *gin.Context) {
	var req dto.LoginRequest
	if result, err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}
	loginResponse, refreshTokenResponse, err := ctrl.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	ctx.SetCookie("refresh_token", refreshTokenResponse.RefreshToken, refreshTokenResponse.ExpiresIn, "/", "", config.GetGlobalConfig().Server.Mode != "debug", true)
	response.Success(ctx, "登录成功", loginResponse)
}

func (ctrl *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		response.Error(ctx, errors.Unauthorized().Msg("刷新令牌为空").Build())
		return
	}
	accessTokenResponse, err := ctrl.authService.RefreshToken(ctx.Request.Context(), refreshToken)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "刷新令牌成功", accessTokenResponse)
}

func (ctrl *AuthController) LogoutHandler(ctx *gin.Context) {
	ctx.SetCookie("refresh_token", "", 0, "/", "", config.GetGlobalConfig().Server.Mode != "debug", true)
	response.Success(ctx, "登出成功", nil)
}