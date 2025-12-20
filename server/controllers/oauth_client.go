package controllers

import (
	"strconv"

	"github.com/3086953492/gokit/ginx/problem"
	"github.com/3086953492/gokit/ginx/response"
	"github.com/3086953492/gokit/validator"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services/oauth"
)

type OAuthClientController struct {
	oauthClientService *oauthservices.OAuthClientService
	validatorManager   *validator.Manager
}

func NewOAuthClientController(oauthClientService *oauthservices.OAuthClientService, validatorManager *validator.Manager) *OAuthClientController {
	return &OAuthClientController{oauthClientService: oauthClientService, validatorManager: validatorManager}
}

func (ctrl *OAuthClientController) CreateOAuthClientHandler(ctx *gin.Context) {
	var req dto.CreateOAuthClientRequest
	if ctx.ShouldBindJSON(&req) != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
		return
	}
	if result := ctrl.validatorManager.Validate(req); !result.Valid {
		problem.Fail(ctx, 400, "INVALID_REQUEST", result.Message, "about:blank")
		return
	}

	if err := ctrl.oauthClientService.CreateOAuthClient(ctx.Request.Context(), &req); err != nil {
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}
	response.OK(ctx, nil, response.WithMessage("创建OAuth客户端成功"))
}

func (ctrl *OAuthClientController) ListOAuthClientsHandler(ctx *gin.Context) {
	page, pageSize := ctx.Query("page"), ctx.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "页码格式错误", "about:blank")
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "每页条数格式错误", "about:blank")
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
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}
	response.OK(ctx, oauthClientsPagination, response.WithMessage("获取OAuth客户端列表成功"))
}

func (ctrl *OAuthClientController) GetOAuthClientHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "ID格式错误", "about:blank")
		return
	}
	oauthClient, err := ctrl.oauthClientService.GetOAuthClient(ctx.Request.Context(), map[string]any{"id": idUint})
	if err != nil {
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}
	response.OK(ctx, oauthClient, response.WithMessage("获取OAuth客户端详情成功"))
}

func (ctrl *OAuthClientController) UpdateOAuthClientHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "ID格式错误", "about:blank")
		return
	}

	var req dto.UpdateOAuthClientRequest
	if ctx.ShouldBindJSON(&req) != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
		return
	}
	if result := ctrl.validatorManager.Validate(req); !result.Valid {
		problem.Fail(ctx, 400, "INVALID_REQUEST", result.Message, "about:blank")
		return
	}

	if err := ctrl.oauthClientService.UpdateOAuthClient(ctx.Request.Context(), uint(idUint), &req); err != nil {
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	response.OK(ctx, nil, response.WithMessage("更新OAuth客户端成功"))
}

func (ctrl *OAuthClientController) DeleteOAuthClientHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		problem.Fail(ctx, 400, "INVALID_REQUEST", "ID格式错误", "about:blank")
		return
	}
	if err := ctrl.oauthClientService.DeleteOAuthClient(ctx.Request.Context(), uint(idUint)); err != nil {
		problem.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	response.OK(ctx, nil, response.WithMessage("删除OAuth客户端成功"))
}
