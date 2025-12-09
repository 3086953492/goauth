package controllers

import (
	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/response"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
	"goauth/utils"
)

type OAuthController struct {
	oauthAuthorizationCodeService *services.OAuthAuthorizationCodeService

	oauthAccessTokenService *services.OAuthAccessTokenService

	oauthClientService *services.OAuthClientService
}

func NewOAuthController(oauthAuthorizationCodeService *services.OAuthAuthorizationCodeService, oauthAccessTokenService *services.OAuthAccessTokenService, oauthClientService *services.OAuthClientService) *OAuthController {
	return &OAuthController{oauthAuthorizationCodeService: oauthAuthorizationCodeService, oauthAccessTokenService: oauthAccessTokenService, oauthClientService: oauthClientService}
}

func (ctrl *OAuthController) AuthorizationCodeHandler(ctx *gin.Context) {

	frontendErrorPageURL := config.GetGlobalConfig().Server.FrontendURL + "/error"

	if responseType := ctx.Query("response_type"); responseType == "" || responseType != "code" {
		response.RedirectTemporary(ctx, frontendErrorPageURL, errors.InvalidInput().Msg("response_type错误").Build(), nil)
		return
	}

	clientID := ctx.Query("client_id")
	if clientID == "" {
		response.RedirectTemporary(ctx, frontendErrorPageURL, errors.InvalidInput().Msg("client_id不能为空").Build(), nil)
		return
	}

	oauthClient, err := ctrl.oauthClientService.GetOAuthClient(ctx.Request.Context(), map[string]any{"id": clientID})
	if err != nil {
		response.RedirectTemporary(ctx, frontendErrorPageURL, err, nil)
		return
	}

	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" || !utils.IsRedirectURIValid(redirectURI, oauthClient.RedirectURIs) {
		response.RedirectTemporary(ctx, frontendErrorPageURL, errors.InvalidInput().Msg("redirect_uri为空或不在客户端的回调地址列表中").Build(), nil)
		return
	}

	state := ctx.Query("state")

	scope := ctx.Query("scope")
	if !utils.IsScopeValid(scope, oauthClient.Scopes) {
		response.RedirectTemporary(ctx, redirectURI, errors.InvalidInput().Msg("scope不在客户端的权限范围列表中").Build(), map[string]string{"state": state})
		return
	}

	userID := uint(ctx.GetUint64("user_id"))

	authorizationCode, err := ctrl.oauthAuthorizationCodeService.GenerateAuthorizationCode(ctx.Request.Context(), userID, clientID, redirectURI, scope)
	if err != nil {
		response.RedirectTemporary(ctx, redirectURI, err, map[string]string{"state": state})
		return
	}

	response.RedirectTemporary(ctx, redirectURI, nil, map[string]string{"state": state, "code": authorizationCode})
}

func (ctrl *OAuthController) ExchangeAccessTokenHandler(ctx *gin.Context) {
	var form dto.ExchangeAccessTokenForm
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", form).Build())
		return
	}

	clientID, clientSecret, ok := ctx.Request.BasicAuth()
	if !ok || clientID == "" || clientSecret == "" {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Build())
		return
	}

	accessToken, err := ctrl.oauthAccessTokenService.ExchangeAccessToken(ctx.Request.Context(), &form, clientID, clientSecret)
	if err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "访问令牌交换成功", accessToken)
}

func (ctrl *OAuthController) IntrospectAccessTokenHandler(ctx *gin.Context) {
	// 绑定请求参数
	var form dto.IntrospectionRequest
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", form).Build())
		return
	}

	// 客户端 Basic 认证
	clientID, clientSecret, ok := ctx.Request.BasicAuth()
	if !ok || clientID == "" || clientSecret == "" {
		response.Error(ctx, errors.InvalidInput().Msg("客户端认证失败").Build())
		return
	}

	// 验证客户端合法性
	_, err := ctrl.oauthClientService.GetOAuthClient(ctx.Request.Context(), map[string]any{"id": clientID, "client_secret": clientSecret})
	if err != nil {
		response.Error(ctx, err)
		return
	}

	// 调用服务层内省访问令牌
	resp := ctrl.oauthAccessTokenService.IntrospectAccessToken(ctx.Request.Context(), form.Token)

	response.Success(ctx, "内省成功", resp)
}
