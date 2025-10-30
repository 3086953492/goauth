package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuthScope struct {
	ID          uint           `gorm:"type:bigint;comment:权限范围ID;primaryKey" json:"id"`
	CreatedAt   time.Time      `gorm:"type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime;comment:更新时间" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;comment:删除时间;index" json:"-"`
	Name        string         `gorm:"type:varchar(50);comment:权限范围名称;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:varchar(255);comment:权限描述" json:"description"`
	IsDefault   bool           `gorm:"type:tinyint(1);comment:是否默认权限;default:false" json:"is_default"`
	Status      int            `gorm:"type:tinyint;comment:状态;default:1" json:"status"` // 1:启用 0:禁用
}

func (OAuthScope) TableName() string {
	return "oauth_scopes"
}

