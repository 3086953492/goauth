package controllers

import (
	"strconv"

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

func (ctrl *OAuthClientController) ListOAuthClientsHandler(ctx *gin.Context) {
	page, pageSize := ctx.Query("page"), ctx.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("页码格式错误").Err(err).Build())
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("每页条数格式错误").Err(err).Build())
		return
	}
	conds := map[string]any{}
	if name := ctx.Query("name"); name != "" {
		conds["name LIKE ?"] = "%" + name + "%"
	}
	if status := ctx.Query("status"); status != "" {
		conds["status"] = status
	}

	oauthClientsPagination, err := ctrl.oauthClientService.ListOAuthClients(ctx.Request.Context(), pageInt, pageSizeInt, conds)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "获取OAuth客户端列表成功", oauthClientsPagination)
}
