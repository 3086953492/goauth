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

type UpdateUserRequest struct { // *字段传递空值会更新为空值，不传递则不更新
	Username string  `json:"username" validate:"omitempty,min=3,max=20,username_unique"`
	Nickname string  `json:"nickname" validate:"omitempty,min=1,max=20"`
	Avatar   *string `json:"avatar" validate:"omitempty,url"`
	Status   *int    `json:"status" validate:"omitempty,oneof=1 0"`
	Role     string  `json:"role" validate:"omitempty,oneof=admin user"`
}
