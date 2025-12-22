package oauthcontrollers

import (
	"strings"

	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/ginx/response"
	"github.com/gin-gonic/gin"

	oauthservices "goauth/services/oauth"
)

type OAuthUserInfoController struct {
	oauthUserInfoService   *oauthservices.OAuthUserInfoService
	OAuthIntrospectService *oauthservices.OAuthIntrospectService
}

func NewOAuthUserInfoController(oauthUserInfoService *oauthservices.OAuthUserInfoService, OAuthIntrospectService *oauthservices.OAuthIntrospectService) *OAuthUserInfoController {
	return &OAuthUserInfoController{oauthUserInfoService: oauthUserInfoService, OAuthIntrospectService: OAuthIntrospectService}
}

func (ctrl *OAuthUserInfoController) GetUserInfoHandler(ctx *gin.Context) {
	// 解析 Bearer Token
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		problem.Fail(ctx, 401, "INVALID_REQUEST", "缺少或无效的 Authorization header", "about:blank")
		return
	}
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// 验证 Token
	introspectResp := ctrl.OAuthIntrospectService.IntrospectAccessToken(ctx.Request.Context(), accessToken)
	if !introspectResp.Active {
		problem.Fail(ctx, 401, "INVALID_TOKEN", "无效的令牌", "about:blank")
		return
	}

	// 校验 scope 是否包含 profile
	if !strings.Contains(introspectResp.Scope, "profile") {
		problem.Fail(ctx, 403, "INSUFFICIENT_SCOPE", "令牌缺少 profile 权限", "about:blank")
		return
	}

	// 获取用户信息
	userInfo := ctrl.oauthUserInfoService.GetUserInfo(ctx.Request.Context(), introspectResp.Username)
	if userInfo == nil {
		problem.Fail(ctx, 404, "USER_NOT_FOUND", "用户不存在", "about:blank")
		return
	}

	response.OK(ctx, userInfo, response.WithMessage("获取用户信息成功"))
}
