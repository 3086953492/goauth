package oauthcontrollers

import (
	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/ginx/response"
	"github.com/gin-gonic/gin"

	"goauth/dto/oauth"
	"goauth/services/oauth"
)

type OAuthIntrospectController struct {
	oauthIntrospectService *oauthservices.OAuthIntrospectService
	oauthClientService *oauthservices.OAuthClientService
}

func NewOAuthIntrospectController(oauthIntrospectService *oauthservices.OAuthIntrospectService, oauthClientService *oauthservices.OAuthClientService) *OAuthIntrospectController {
	return &OAuthIntrospectController{oauthIntrospectService: oauthIntrospectService, oauthClientService: oauthClientService}
}

func (ctrl *OAuthIntrospectController) IntrospectAccessTokenHandler(ctx *gin.Context) {
	// 绑定请求参数
	var form oauthdto.IntrospectionRequest
	if err := ctx.ShouldBind(&form); err != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
		return
	}

	// 客户端 Basic 认证
	clientID, clientSecret, ok := ctx.Request.BasicAuth()
	if !ok || clientID == "" || clientSecret == "" {
		problem.Fail(ctx, 401, "INVALID_CLIENT", "非法的客户端凭证", "about:blank")
		return
	}

	// 验证客户端合法性
	_, err := ctrl.oauthClientService.GetOAuthClient(ctx.Request.Context(), map[string]any{"id": clientID, "client_secret": clientSecret})
	if err != nil {
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	// 调用服务层内省访问令牌
	resp := ctrl.oauthIntrospectService.IntrospectAccessToken(ctx.Request.Context(), form.Token)

	response.OK(ctx, resp, response.WithMessage("内省访问令牌成功"))
}