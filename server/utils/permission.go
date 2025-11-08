package utils

import (
	"github.com/gin-gonic/gin"
)

func IsRole(c *gin.Context, requiredRole string) bool {
	if requiredRole == "" {
		return true
	}
	if role, ok := c.Get("role"); !ok || role != requiredRole {
		return false
	}
	return true
}

func IsResourceOwner(c *gin.Context, source string) bool {
	if source == "" || IsRole(c, "admin") {
		return true
	}

	var requiredUserID string
	switch source {
	case "param":
		requiredUserID = c.Param("user_id")
	case "query":
		requiredUserID = c.Query("user_id")
	default:
		return false
	}
	
	if requiredUserID == "" {
		return false
	}
	currentUserID, ok := c.Get("user_id")
	if !ok {
		return false
	}
	currentUserIDStr, ok := currentUserID.(string)
	if !ok || currentUserIDStr == "" {
		return false
	}
	if currentUserIDStr != requiredUserID {
		return false
	}
	return true
}
