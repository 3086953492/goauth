package controllers

import (
	"strconv"

	"github.com/3086953492/gokit/ginx"
	"github.com/3086953492/gokit/validator"
	"github.com/gin-gonic/gin"

	"goauth/dto"
	"goauth/services"
	"goauth/utils"
)

type UserController struct {
	userService      *services.UserService
	validatorManager *validator.Manager
}

func NewUserController(userService *services.UserService, validatorManager *validator.Manager) *UserController {
	return &UserController{userService: userService, validatorManager: validatorManager}
}

func (ctrl *UserController) CreateUserHandler(ctx *gin.Context) {
	var form dto.CreateUserForm
	if ctx.ShouldBind(&form) != nil {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
		return
	}
	if result := ctrl.validatorManager.Validate(form); !result.Valid {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", result.Message, "about:blank")
		return
	}

	// 校验头像文件（可选），不打开文件
	avatarFile, err := utils.ValidateFormFile(ctx, "avatar", 4*1024*1024, []string{"image/png", "image/jpeg", "image/jpg", "image/webp"})
	if err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	if err := ctrl.userService.CreateUser(ctx.Request.Context(), &form, avatarFile); err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	ginx.OK(ctx, nil)
}

func (ctrl *UserController) UpdateUserHandler(ctx *gin.Context) {
	var form dto.UpdateUserForm
	if ctx.ShouldBind(&form) != nil {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", "请求参数错误", "about:blank")
		return
	}
	if result := ctrl.validatorManager.Validate(form); !result.Valid {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", result.Message, "about:blank")
		return
	}

	if !utils.IsRole(ctx, "admin") { // 非管理员不能修改状态和角色
		form.Status = nil
		form.Role = ""
	}

	userID := ctx.Param("user_id")
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", "用户ID格式错误", "about:blank")
		return
	}

	avatarFile, err := utils.ValidateFormFile(ctx, "avatar", 4*1024*1024, []string{"image/png", "image/jpeg", "image/jpg", "image/webp"})
	if err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}

	if err := ctrl.userService.UpdateUser(ctx.Request.Context(), uint(userIDUint), &form, avatarFile); err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}
	ginx.OK(ctx, nil)
}

func (ctrl *UserController) GetUserHandler(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	user, err := ctrl.userService.GetUser(ctx.Request.Context(), map[string]any{"id": userID})
	if err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}
	ginx.OK(ctx, user)
}

func (ctrl *UserController) ListUsersHandler(ctx *gin.Context) {
	page, pageSize := ctx.Query("page"), ctx.Query("page_size")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", "页码格式错误", "about:blank")
		return
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", "每页条数格式错误", "about:blank")
		return
	}

	status := ctx.Query("status")
	role := ctx.Query("role")
	nickname := ctx.Query("nickname")
	conds := map[string]any{}
	if status != "" {
		conds["status"] = status
	}
	if role != "" {
		conds["role"] = role
	}
	if nickname != "" {
		conds["nickname LIKE ?"] = "%" + nickname + "%"
	}

	users, err := ctrl.userService.ListUsers(ctx.Request.Context(), pageInt, pageSizeInt, conds)
	if err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}
	ginx.OK(ctx, users)
}

func (ctrl *UserController) DeleteUserHandler(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		ginx.Fail(ctx, 400, "INVALID_REQUEST", "用户ID格式错误", "about:blank")
		return
	}

	if uint(userIDUint) == ctx.GetUint("user_id") {
		ginx.Fail(ctx, 403, "FORBIDDEN", "不能删除自己", "about:blank")
		return
	}

	if err := ctrl.userService.DeleteUser(ctx.Request.Context(), uint(userIDUint)); err != nil {
		ginx.Fail(ctx, 500, "INTERNAL_SERVER_ERROR", err.Error(), "about:blank")
		return
	}
	ginx.OK(ctx, nil)
}
