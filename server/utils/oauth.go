package utils

import (
	"encoding/json"
	"slices"
	"strings"
)

// isRedirectURIValid 验证 redirect_uri 是否在白名单中（精确匹配）
func IsRedirectURIValid(redirectURI string, registeredURIsJSON []byte) bool {
	var registeredURIs []string
	if err := json.Unmarshal(registeredURIsJSON, &registeredURIs); err != nil {
		return false
	}

	return slices.Contains(registeredURIs, redirectURI)
}

// isScopeValid 验证请求的 scope 是否都在允许的 scope 列表中
func IsScopeValid(requestedScope string, allowedScopesJSON []byte) bool {
	if requestedScope == "" {
		return true // 空 scope 视为合法
	}

	var allowedScopes []string
	if err := json.Unmarshal(allowedScopesJSON, &allowedScopes); err != nil {
		return false
	}

	// 将允许的 scope 转为 map 便于查找
	allowedScopeMap := make(map[string]bool)
	for _, scope := range allowedScopes {
		allowedScopeMap[scope] = true
	}

	// 检查请求的每个 scope 是否都在允许列表中
	requestedScopes := strings.Split(requestedScope, " ")
	for _, scope := range requestedScopes {
		scope = strings.TrimSpace(scope)
		if scope != "" && !allowedScopeMap[scope] {
			return false // 有任何一个 scope 不在允许列表中就返回 false
		}
	}

	return true
}
