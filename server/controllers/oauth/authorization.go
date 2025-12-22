package oauthcontrollers

import (
	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/ginx/redirect"
	"github.com/gin-gonic/gin"

	"goauth/services/oauth"
	"goauth/utils"
)

type OAuthAuthorizationController struct {
	oauthAuthorizationCodeService *oauthservices.OAuthAuthorizationCodeService

	oauthAccessTokenService *oauthservices.OAuthAccessTokenService

	oauthClientService *oauthservices.OAuthClientService

	cfg *config.Config
}

func NewOAuthAuthorizationController(oauthAuthorizationCodeService *oauthservices.OAuthAuthorizationCodeService, oauthAccessTokenService *oauthservices.OAuthAccessTokenService, oauthClientService *oauthservices.OAuthClientService, cfg *config.Config) *OAuthAuthorizationController {
	return &OAuthAuthorizationController{oauthAuthorizationCodeService: oauthAuthorizationCodeService, oauthAccessTokenService: oauthAccessTokenService, oauthClientService: oauthClientService, cfg: cfg}
}

func (ctrl *OAuthAuthorizationController) AuthorizationCodeHandler(ctx *gin.Context) {

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
