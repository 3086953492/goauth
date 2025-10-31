package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Nickname string `json:"nickname" validate:"required,min=1,max=20"`
	Avatar   string `json:"avatar" validate:"omitempty,url"`
}
