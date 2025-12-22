package oauthcontrollers

import (
	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/ginx/response"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services/oauth"
)

type OAuthTokenController struct {
	oauthAccessTokenService *oauthservices.OAuthAccessTokenService

	oauthClientService *oauthservices.OAuthClientService
}

func NewOAuthTokenController(oauthAccessTokenService *oauthservices.OAuthAccessTokenService, oauthClientService *oauthservices.OAuthClientService) *OAuthTokenController {
	return &OAuthTokenController{oauthAccessTokenService: oauthAccessTokenService, oauthClientService: oauthClientService}
}

func (ctrl *OAuthTokenController) ExchangeAccessTokenHandler(ctx *gin.Context) {
	// 客户端 Basic 认证
	clientID, clientSecret, ok := ctx.Request.BasicAuth()
	if !ok || clientID == "" || clientSecret == "" {
		problem.Fail(ctx, 401, "INVALID_CLIENT", "非法的客户端凭证", "about:blank")
		return
	}

	// 根据 grant_type 分支处理
	grantType := ctx.PostForm("grant_type")
	switch grantType {
	case "authorization_code":
		var form dto.ExchangeAccessTokenForm
		if err := ctx.ShouldBind(&form); err != nil {
			problem.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
			return
		}

		accessToken, err := ctrl.oauthAccessTokenService.ExchangeAccessToken(ctx.Request.Context(), &form, clientID, clientSecret)
		if err != nil {
			problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
			return
		}

		response.OK(ctx, accessToken, response.WithMessage("交换访问令牌成功"))

	case "refresh_token":
		var form dto.RefreshAccessTokenForm
		if err := ctx.ShouldBind(&form); err != nil {
			problem.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
			return
		}

		accessToken, err := ctrl.oauthAccessTokenService.RefreshAccessToken(ctx.Request.Context(), &form, clientID, clientSecret)
		if err != nil {
			problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
			return
		}

		response.OK(ctx, accessToken, response.WithMessage("刷新访问令牌成功"))

	default:
		problem.Fail(ctx, 400, "INVALID_REQUEST", "授权类型不支持", "about:blank")
	}
}



func (ctrl *OAuthTokenController) IntrospectAccessTokenHandler(ctx *gin.Context) {
	// 绑定请求参数
	var form dto.IntrospectionRequest
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
	resp := ctrl.oauthAccessTokenService.IntrospectAccessToken(ctx.Request.Context(), form.Token)

	response.OK(ctx, resp, response.WithMessage("内省访问令牌成功"))
}
