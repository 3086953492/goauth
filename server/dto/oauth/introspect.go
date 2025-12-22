package oauthdto

// IntrospectionRequest 内省请求结构体
type IntrospectionRequest struct {
	Token         string `form:"token" binding:"required"`
	TokenTypeHint string `form:"token_type_hint"` // "access_token", "refresh_token"
}

// IntrospectionResponse 内省响应结构体 (RFC 7662)
type IntrospectionResponse struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	ClientID  string `json:"client_id,omitempty"`
	Username  string `json:"username,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Exp       int64  `json:"exp,omitempty"`
	Sub       string `json:"sub,omitempty"`
}
