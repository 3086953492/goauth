package oauthdto

type UserInfoResponse struct {
	Sub       string `json:"sub"`
	Nickname  string `json:"nickname"`
	Picture   string `json:"picture"`
	UpdatedAt int64  `json:"updated_at"`
}
