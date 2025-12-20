package controllers

import (
	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/ginx/redirect"
	"github.com/3086953492/gokit/ginx/response"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services/oauth"
	"goauth/utils"
)

type OAuthController struct {
	oauthAuthorizationCodeService *oauthservices.OAuthAuthorizationCodeService

	oauthAccessTokenService *oauthservices.OAuthAccessTokenService

	oauthClientService *oauthservices.OAuthClientService

	cfg *config.Config
}

func NewOAuthController(oauthAuthorizationCodeService *oauthservices.OAuthAuthorizationCodeService, oauthAccessTokenService *oauthservices.OAuthAccessTokenService, oauthClientService *oauthservices.OAuthClientService, cfg *config.Config) *OAuthController {
	return &OAuthController{oauthAuthorizationCodeService: oauthAuthorizationCodeService, oauthAccessTokenService: oauthAccessTokenService, oauthClientService: oauthClientService, cfg: cfg}
}

func (ctrl *OAuthController) AuthorizationCodeHandler(ctx *gin.Context) {

	frontendErrorPageURL := ctrl.cfg.Server.FrontendURL + "/error"

	if responseType := ctx.Query("response_type"); responseType == "" || responseType != "code" {
		redirect.Redirect(ctx, frontendErrorPageURL, redirect.WithQuery(map[string]string{"error": "invalid_request", "error_description": "response_type错误"}))
		return
	}

	clientID := ctx.Query("client_id")
	if clientID == "" {
		redirect.Redirect(ctx, frontendErrorPageURL, redirect.WithQuery(map[string]string{"error": "invalid_request", "error_description": "client_id不能为空"}))
		return
	}

	oauthClient, err := ctrl.oauthClientService.GetOAuthClient(ctx.Request.Context(), map[string]any{"id": clientID})
	if err != nil {
		redirect.Redirect(ctx, frontendErrorPageURL, redirect.WithQuery(map[string]string{"error": "invalid_request", "error_description": err.Error()}))
		return
	}

	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" || !utils.IsRedirectURIValid(redirectURI, oauthClient.RedirectURIs) {
		redirect.Redirect(ctx, frontendErrorPageURL, redirect.WithQuery(map[string]string{"error": "invalid_request", "error_description": "redirect_uri为空或不在客户端的回调地址列表中"}))
		return
	}

	state := ctx.Query("state")

	scope := ctx.Query("scope")
	if !utils.IsScopeValid(scope, oauthClient.Scopes) {
		redirect.Redirect(ctx, redirectURI, redirect.WithQuery(map[string]string{"error": "invalid_request", "error_description": "scope不在客户端的权限范围列表中"}))
		return
	}

	userID := uint(ctx.GetUint64("user_id"))

	authorizationCode, err := ctrl.oauthAuthorizationCodeService.GenerateAuthorizationCode(ctx.Request.Context(), userID, clientID, redirectURI, scope)
	if err != nil {
		redirect.Redirect(ctx, redirectURI, redirect.WithQuery(map[string]string{"error": "invalid_request", "error_description": err.Error()}))
		return
	}

	redirect.Redirect(ctx, redirectURI, redirect.WithQuery(map[string]string{"state": state, "code": authorizationCode}))
}

func (ctrl *OAuthController) ExchangeAccessTokenHandler(ctx *gin.Context) {
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

func (ctrl *OAuthController) IntrospectAccessTokenHandler(ctx *gin.Context) {
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
