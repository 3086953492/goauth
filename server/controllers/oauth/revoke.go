package oauthcontrollers

import (
	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/ginx/response"
	"github.com/gin-gonic/gin"

	"goauth/dto/oauth"
	"goauth/services/oauth"
)

// OAuthRevokeController 令牌撤销控制器（RFC7009）
type OAuthRevokeController struct {
	oauthRevokeService *oauthservices.OAuthRevokeService
	oauthClientService *oauthservices.OAuthClientService
}

// NewOAuthRevokeController 创建令牌撤销控制器实例
func NewOAuthRevokeController(oauthRevokeService *oauthservices.OAuthRevokeService, oauthClientService *oauthservices.OAuthClientService) *OAuthRevokeController {
	return &OAuthRevokeController{
		oauthRevokeService: oauthRevokeService,
		oauthClientService: oauthClientService,
	}
}

// RevokeTokenHandler 处理令牌撤销请求
func (ctrl *OAuthRevokeController) RevokeTokenHandler(ctx *gin.Context) {
	// 客户端 Basic 认证
	clientID, clientSecret, ok := ctx.Request.BasicAuth()
	if !ok || clientID == "" || clientSecret == "" {
		problem.Fail(ctx, 401, "INVALID_CLIENT", "非法的客户端凭证", "about:blank")
		return
	}

	// 验证客户端合法性
	_, err := ctrl.oauthClientService.GetOAuthClient(ctx.Request.Context(), map[string]any{"id": clientID, "client_secret": clientSecret})
	if err != nil {
		problem.Fail(ctx, 401, "INVALID_CLIENT", "非法的客户端凭证", "about:blank")
		return
	}

	// 绑定请求参数
	var form oauthdto.RevocationRequest
	if err := ctx.ShouldBind(&form); err != nil {
		// 检查是否是 token_type_hint 校验失败
		if form.Token == "" {
			problem.Fail(ctx, 400, "INVALID_REQUEST", "缺少必填参数 token", "about:blank")
			return
		}
		// token_type_hint 非法值
		problem.Fail(ctx, 400, "UNSUPPORTED_TOKEN_TYPE", "不支持的 token_type_hint 值", "about:blank")
		return
	}

	// 调用服务层撤销令牌
	_ = ctrl.oauthRevokeService.RevokeToken(ctx.Request.Context(), form.Token, form.TokenTypeHint, clientID)

	// RFC7009：无论是否成功撤销，均返回 200
	response.OK(ctx, nil, response.WithMessage("令牌撤销成功"))
}

