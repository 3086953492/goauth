package controllers

import (
	"github.com/3086953492/gokit/errors"
	"github.com/3086953492/gokit/response"
	"github.com/3086953492/gokit/validator"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
	"goauth/utils"
)

type UserController struct {
	userService *services.UserService
}


func (ctrl *UserController) CreateUserHandler(ctx *gin.Context) {
	var req dto.CreateUserRequest
	if result, err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}

	if err := ctrl.userService.CreateUser(ctx.Request.Context(), &req); err != nil {
		response.Error(ctx, err)
		return
	}

	response.Success(ctx, "创建用户成功", nil)
}

func (ctrl *UserController) UpdateUserHandler(ctx *gin.Context) {
	userID := ctx.Param("id")
	var req dto.UpdateUserRequest
	if result, err := validator.BindAndValidate(ctx, &req); err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("请求参数错误").Err(err).Field("request", req).Build())
		return
	} else if !result.Valid {
		response.Error(ctx, errors.InvalidInput().Msg(result.Message).Err(result.Err).Field("request", req).Build())
		return
	}

	if !utils.IsRole(ctx, "admin") { // 非管理员不能修改状态和角色
		req.Status = nil
		req.Role = ""
	}

	if err := ctrl.userService.UpdateUser(ctx.Request.Context(), &dto.UpdateUserRequest{
		Username: userID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Role:     req.Role,
	}); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "更新用户成功", nil)
}

func (ctrl *UserController) GetUserHandler(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	user, err := ctrl.userService.GetUser(ctx.Request.Context(), map[string]any{"id": userID})
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "获取用户成功", user)
}