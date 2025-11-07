package dto

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=20,username_unique"`
	Password        string `json:"password" validate:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Nickname        string `json:"nickname" validate:"required,min=1,max=20"`
	Avatar          string `json:"avatar" validate:"omitempty,url"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginResponse struct {
	User  *UserResponse  `json:"user"`
	Token *TokenResponse `json:"token"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
