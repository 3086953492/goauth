package controllers

import (
	"time"

	// "github.com/3086953492/gokit/cookie"
	"github.com/3086953492/gokit/ginx/cookie"
	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/ginx/response"
	"github.com/3086953492/gokit/validator"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
)

type AuthController struct {
	authService      *services.AuthService
	validatorManager *validator.Manager
	cookieMgr *cookie.TokenCookies
}

func NewAuthController(authService *services.AuthService, validatorManager *validator.Manager, cookieMgr *cookie.TokenCookies) *AuthController {
	return &AuthController{authService: authService, validatorManager: validatorManager, cookieMgr: cookieMgr}
}

func (ctrl *AuthController) LoginHandler(ctx *gin.Context) {
	var req dto.LoginRequest
	if ctx.ShouldBindJSON(&req) != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
		return
	}
	if result := ctrl.validatorManager.Validate(req); !result.Valid {
		problem.Fail(ctx, 400, "INVALID_REQUEST", result.Message, "about:blank")
		return
	}

	accessToken, accessTokenExpire, refreshToken, refreshTokenExpire, userResp, err := ctrl.authService.Login(ctx.Request.Context(), &req)
	if err != nil {
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	ctrl.cookieMgr.SetAccess(ctx, accessToken)
	ctrl.cookieMgr.SetRefresh(ctx, refreshToken)

	response.OK(ctx, dto.LoginResponse{
		User:                 userResp,
		AccessTokenExpireAt:  time.Now().Add(time.Duration(accessTokenExpire) * time.Second),
		RefreshTokenExpireAt: time.Now().Add(time.Duration(refreshTokenExpire) * time.Second),
	},response.WithMessage("登录成功"))
}

func (ctrl *AuthController) LogoutHandler(ctx *gin.Context) {
	ctrl.cookieMgr.Clear(ctx)
	response.OK(ctx, nil, response.WithMessage("退出登录成功"))
}

func (ctrl *AuthController) RefreshTokenHandler(ctx *gin.Context) {
	token, err := ctrl.cookieMgr.GetRefresh(ctx)
	if err != nil {
		problem.Fail(ctx, 401, "UNAUTHORIZED", "刷新令牌不存在", "about:blank")
		return
	}
	accessToken, accessTokenExpire, err := ctrl.authService.RefreshToken(ctx.Request.Context(), token)
	if err != nil {
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	ctrl.cookieMgr.SetAccess(ctx, accessToken)

	response.OK(ctx, dto.RefreshTokenResponse{
		AccessTokenExpireAt: time.Now().Add(time.Duration(accessTokenExpire) * time.Second),
	}, response.WithMessage("刷新令牌成功"))
}
