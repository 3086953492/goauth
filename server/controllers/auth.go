package controllers

import (
	"time"

	"github.com/3086953492/gokit/cookie"
	"github.com/3086953492/gokit/ginx"
	"github.com/3086953492/gokit/validator"
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
	if ctx.ShouldBindJSON(&req) != nil {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
		return
	}
	if result := ctrl.validatorManager.Validate(req); !result.Valid {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", result.Message, "about:blank")
		return
	}

	accessToken, accessTokenExpire, refreshToken, refreshTokenExpire, userResp, err := ctrl.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	cookie.SetAccessToken(ctx, accessToken, accessTokenExpire, "", "/")
	cookie.SetRefreshToken(ctx, refreshToken, refreshTokenExpire, "", "/")

	ginx.OK(ctx, dto.LoginResponse{
		User:                 userResp,
		AccessTokenExpireAt:  time.Now().Add(time.Duration(accessTokenExpire) * time.Second),
		RefreshTokenExpireAt: time.Now().Add(time.Duration(refreshTokenExpire) * time.Second),
	})
}

func (ctrl *AuthController) LogoutHandler(ctx *gin.Context) {
	cookie.ClearTokens(ctx, "", "/")
	ginx.OK(ctx, nil)
}

func (ctrl *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	token, err := cookie.GetRefreshToken(ctx)
	if err != nil {
		ginx.Fail(ctx, 401, "UNAUTHORIZED", "刷新令牌不存在", "about:blank")
		return
	}
	accessToken, accessTokenExpire, err := ctrl.authService.RefreshToken(ctx.Request.Context(), token)
	if err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	cookie.SetAccessToken(ctx, accessToken, accessTokenExpire, "", "/")

	ginx.OK(ctx, dto.RefreshTokenResponse{
		AccessTokenExpireAt: time.Now().Add(time.Duration(accessTokenExpire) * time.Second),
	})
}
