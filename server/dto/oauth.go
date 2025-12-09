package dto

type AuthorizationCodeResponse struct {
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state"`
}

type ExchangeAccessTokenResponse struct {
	AccessToken  OAuthAccessTokenResponse  `json:"access_token"`
	RefreshToken OAuthRefreshTokenResponse `json:"refresh_token"`
	TokenType    string                    `json:"token_type"`
	Scope        string                    `json:"scope"`
}

type OAuthAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type OAuthRefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type ExchangeAccessTokenForm struct {
	GrantType   string `form:"grant_type" binding:"required,oneof=authorization_code"`
	Code        string `form:"code" binding:"required"`
	RedirectURI string `form:"redirect_uri" binding:"required"`
}

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
