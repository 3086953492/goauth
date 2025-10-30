package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"type:bigint;comment:用户ID;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime;comment:更新时间" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime;comment:删除时间;index" json:"-"`
	Username  string         `gorm:"type:varchar(50);comment:用户名;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);comment:密码哈希;not null" json:"-"`
	Nickname  string         `gorm:"type:varchar(100);comment:昵称" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(500);comment:头像URL" json:"avatar"`
	Status    int            `gorm:"type:tinyint;comment:状态;default:1" json:"status"` // 1:正常 0:禁用
	Role      string         `gorm:"type:varchar(50);comment:角色" json:"role"`
}

func (User) TableName() string {
	return "users"
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
}