package dto

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Status    int       `json:"status"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserForm struct { // *字段传递空值会更新为空值，不传递则不更新
	Nickname        string `form:"nickname" validate:"omitempty,min=1,max=20"`
	Password        string `form:"password" validate:"omitempty,min=6,max=20"`
	ConfirmPassword string `form:"confirm_password" validate:"omitempty,eqfield=Password"`
	Status          *int   `form:"status" validate:"omitempty,oneof=1 0"`
	Role            string `form:"role" validate:"omitempty,oneof=admin user"`
}

type CreateUserForm struct {
	Username        string `form:"username" validate:"required,min=3,max=20,username_unique"`
	Password        string `form:"password" validate:"required,min=6,max=20"`
	ConfirmPassword string `form:"confirm_password" validate:"required,eqfield=Password"`
	Nickname        string `form:"nickname" validate:"required,min=1,max=20"`
}

type UserListResponse struct {
	ID        uint      `json:"id"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Status    int       `json:"status"`
	Role      string    `json:"role"`
}
