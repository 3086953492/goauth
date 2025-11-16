package dto

type AuthorizationCodeResponse struct {
	Code        string `json:"code"`
	ExpiresIn   int    `json:"expires_in"`
	RedirectURI string `json:"redirect_uri"`
	State       string `json:"state"`
}
