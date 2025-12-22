package oauthdto

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

type RefreshAccessTokenForm struct {
	GrantType    string `form:"grant_type" binding:"required,oneof=refresh_token"`
	RefreshToken string `form:"refresh_token" binding:"required"`
}
