package dto

import "time"

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type LoginResponse struct {
	User                 *UserResponse `json:"user"`
	AccessTokenExpireAt  time.Time     `json:"access_token_expire_at"`
	RefreshTokenExpireAt time.Time     `json:"refresh_token_expire_at"`
}

type RefreshTokenResponse struct {
	AccessTokenExpireAt time.Time `json:"access_token_expire_at"`
}