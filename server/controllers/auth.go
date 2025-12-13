package controllers

import (
	"time"

	"github.com/3086953492/gokit/cookie"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/response"
	"github.com/3086953492/gokit/validator"
	vgin "github.com/3086953492/gokit/validator/provider_gin"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
)

type AuthController struct {
	authService      *services.AuthService
	validatorManager *validator.Manager
}

func NewAuthController(authService *services.AuthService, validatorManager *validator.Manager) *AuthController {
	return &AuthController{authService: authService, validatorManager: validatorManager}
}

func (ctrl *AuthController) LoginHandler(ctx *gin.Context) {
	var req dto.LoginRequest
	if result, err := vgin.BindAndValidate(ctrl.validatorManager, ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}

	accessToken, accessTokenExpire, refreshToken, refreshTokenExpire, userResp, err := ctrl.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	cookie.SetAccessToken(ctx, accessToken, accessTokenExpire, "", "/")
	cookie.SetRefreshToken(ctx, refreshToken, refreshTokenExpire, "", "/")

	response.Success(ctx, "登录成功", dto.LoginResponse{
		User:                 userResp,
		AccessTokenExpireAt:  time.Now().Add(time.Duration(accessTokenExpire) * time.Second),
		RefreshTokenExpireAt: time.Now().Add(time.Duration(refreshTokenExpire) * time.Second),
	})
}

func (ctrl *AuthController) LogoutHandler(ctx *gin.Context) {
	cookie.ClearTokens(ctx, "", "/")
	response.Success(ctx, "登出成功", nil)
}

func (ctrl *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	token, err := cookie.GetRefreshToken(ctx)
	if err != nil {
		response.Error(ctx, errors.Unauthorized().Msg("刷新令牌不存在").Err(err).Build())
		return
	}
	accessToken, accessTokenExpire, err := ctrl.authService.RefreshToken(ctx.Request.Context(), token)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	cookie.SetAccessToken(ctx, accessToken, accessTokenExpire, "", "/")

	response.Success(ctx, "刷新令牌成功", dto.RefreshTokenResponse{
		AccessTokenExpireAt: time.Now().Add(time.Duration(accessTokenExpire) * time.Second),
	})
}
