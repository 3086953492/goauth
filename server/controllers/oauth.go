package controllers

import (
	"github.com/3086953492/gokit/config"
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/response"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
)

type OAuthController struct {
	oauthAuthorizationCodeService *services.OAuthAuthorizationCodeService

	oauthAccessTokenService *services.OAuthAccessTokenService
}

func NewOAuthController(oauthAuthorizationCodeService *services.OAuthAuthorizationCodeService, oauthAccessTokenService *services.OAuthAccessTokenService) *OAuthController {
	return &OAuthController{oauthAuthorizationCodeService: oauthAuthorizationCodeService, oauthAccessTokenService: oauthAccessTokenService}
}

func (ctrl *OAuthController) AuthorizationCodeHandler(ctx *gin.Context) {

	clientID := ctx.Query("client_id")
	if clientID == "" {
		response.Error(ctx, errors.InvalidInput().Msg("client_id不能为空").Build())
		return
	}

	redirectURI := ctx.Query("redirect_uri")
	if redirectURI == "" {
		response.Error(ctx, errors.InvalidInput().Msg("redirect_uri不能为空").Build())
		return
	}

	scope := ctx.Query("scope")

	if responseType := ctx.Query("response_type"); responseType == "" || responseType != "code" {
		response.Error(ctx, errors.InvalidInput().Msg("response_type错误").Build())
		return
	}

	state := ctx.Query("state")

	userID := uint(ctx.GetUint64("user_id"))

	expiresIn := config.GetGlobalConfig().JWT.Expire.Seconds()

	authorizationCode, err := ctrl.oauthAuthorizationCodeService.GenerateAuthorizationCode(ctx.Request.Context(), userID, clientID, redirectURI, scope, expiresIn)
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
