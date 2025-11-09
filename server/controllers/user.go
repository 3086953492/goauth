package controllers

import (
	"strconv"

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

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
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

	userID := ctx.Param("user_id")
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("用户ID格式错误").Err(err).Build())
		return
	}

	if err := ctrl.userService.UpdateUser(ctx.Request.Context(), uint(userIDUint), &dto.UpdateUserRequest{
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Status:   req.Status,
		Role:     req.Role,
		Password: req.Password,
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

func (ctrl *UserController) ListUsersHandler(ctx *gin.Context) {
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
	
	status := ctx.Query("status")
	role := ctx.Query("role")
	conds := map[string]any{}
	if status != "" {
		conds["status"] = status
	}
	if role != "" {
		conds["role"] = role
	}

	users, err := ctrl.userService.ListUsers(ctx.Request.Context(), pageInt, pageSizeInt, conds)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "获取用户列表成功", users)
}

func (ctrl *UserController) DeleteUserHandler(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		response.Error(ctx, errors.InvalidInput().Msg("用户ID格式错误").Err(err).Build())
		return
	}

	if uint(userIDUint) == ctx.GetUint("user_id") {
		response.Error(ctx, errors.Forbidden().Msg("不能删除自己").Build())
		return
	}

	if err := ctrl.userService.DeleteUser(ctx.Request.Context(), uint(userIDUint)); err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, "删除用户成功", nil)
}