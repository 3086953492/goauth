package oauthdto

// RevocationRequest RFC7009 令牌撤销请求
type RevocationRequest struct {
	Token         string `form:"token" binding:"required"`
	TokenTypeHint string `form:"token_type_hint" binding:"omitempty,oneof=access_token refresh_token"`
}

