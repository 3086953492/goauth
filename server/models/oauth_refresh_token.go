package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuthRefreshToken struct {
	ID              uint           `gorm:"type:bigint;comment:刷新令牌ID;primaryKey" json:"id"`
	CreatedAt       time.Time      `gorm:"type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:datetime;comment:更新时间" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"type:datetime;comment:删除时间;index" json:"-"`
	RefreshToken    string         `gorm:"type:varchar(255);comment:刷新令牌;uniqueIndex;not null" json:"refresh_token"`
	AccessTokenID   uint           `gorm:"type:bigint;comment:关联的访问令牌ID;index" json:"access_token_id"`
	ClientID        string         `gorm:"type:varchar(100);comment:客户端ID;index;not null" json:"client_id"`
	UserID          uint           `gorm:"type:bigint;comment:用户ID;index;not null" json:"user_id"`
	Scope           string         `gorm:"type:varchar(500);comment:权限范围" json:"scope"`
	ExpiresAt       time.Time      `gorm:"type:datetime;comment:过期时间;index;not null" json:"expires_at"`
	Revoked         bool           `gorm:"type:tinyint(1);comment:是否已撤销;default:false" json:"revoked"`
}

func (OAuthRefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

