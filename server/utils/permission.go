package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func IsRole(c *gin.Context, requiredRole string) bool {
	if requiredRole == "" {
		return true
	}
	roleVal, ok := c.Get("role")
	if !ok {
		return false
	}
	role, ok := roleVal.(string)
	if !ok || role != requiredRole {
		return false
	}
	return true
}

// IsResourceOwner 检查当前用户是否是资源所有者
// source: "param" 从路径参数获取 user_id，"query" 从查询参数获取 user_id
// 注意：
// - 管理员（role=admin）默认通过
// - 客户端主体（principal_kind=client）不受资源所有者限制
// - user_id 在 context 中存储为 uint64 类型
func IsResourceOwner(c *gin.Context, source string) bool {
	if source == "" || IsRole(c, "admin") {
		return true
	}

	// 检查是否为客户端主体（client_credentials 模式）
	// 客户端主体不受资源所有者限制
	if kindVal, exists := c.Get("principal_kind"); exists {
		if kind, ok := kindVal.(string); ok && kind == "client" {
			return true
		}
	}

	// 获取请求中的目标 user_id（字符串形式）
	var requiredUserIDStr string
	switch source {
	case "param":
		requiredUserIDStr = c.Param("user_id")
	case "query":
		requiredUserIDStr = c.Query("user_id")
	default:
		return false
	}

	if requiredUserIDStr == "" {
		return false
	}

	// 将请求中的 user_id 解析为数字
	requiredUserID, err := strconv.ParseUint(requiredUserIDStr, 10, 64)
	if err != nil {
		return false
	}

	// 获取当前用户 ID（uint64 类型）
	currentUserIDVal, ok := c.Get("user_id")
	if !ok {
		return false
	}

	currentUserID, ok := currentUserIDVal.(uint64)
	if !ok || currentUserID == 0 {
		return false
	}

	return currentUserID == requiredUserID
}
