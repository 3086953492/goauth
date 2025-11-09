package controllers

import (
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/response"
	"github.com/3086953492/gokit/validator"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
)

type OAuthClientController struct {
	oauthClientService *services.OAuthClientService
}

func NewOAuthClientController(oauthClientService *services.OAuthClientService) *OAuthClientController {
	return &OAuthClientController{oauthClientService: oauthClientService}
}

func (ctrl *OAuthClientController) CreateOAuthClientHandler(ctx *gin.Context) {
	var req dto.CreateOAuthClientRequest
	if result, err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}
	if err := ctrl.oauthClientService.CreateOAuthClient(ctx.Request.Context(), &req); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "创建OAuth客户端成功", nil)
}