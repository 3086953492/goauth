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

func (ctrl *OAuthClientController) GetOAuthClientHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("ID格式错误").Err(err).Build())
		return
	}
	oauthClient, err := ctrl.oauthClientService.GetOAuthClient(ctx.Request.Context(), map[string]any{"id": idUint})
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "获取OAuth客户端成功", oauthClient)
}

func (ctrl *OAuthClientController) UpdateOAuthClientHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("ID格式错误").Err(err).Build())
		return
	}

	var req dto.UpdateOAuthClientRequest
	if result, err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}

	if err := ctrl.oauthClientService.UpdateOAuthClient(ctx.Request.Context(), uint(idUint), &req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "更新OAuth客户端成功", nil)
}

func (ctrl *OAuthClientController) DeleteOAuthClientHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("ID格式错误").Err(err).Build())
		return
	}
	if err := ctrl.oauthClientService.DeleteOAuthClient(ctx.Request.Context(), uint(idUint)); err != nil {
		response.Error(ctx, err)
		return
	}
	
	response.Success(ctx, "删除OAuth客户端成功", nil)
}