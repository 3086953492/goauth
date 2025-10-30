package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuthAccessToken struct {
	ID          uint           `gorm:"type:bigint;comment:令牌ID;primaryKey" json:"id"`
	CreatedAt   time.Time      `gorm:"type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime;comment:更新时间" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;comment:删除时间;index" json:"-"`
	AccessToken string         `gorm:"type:varchar(255);comment:访问令牌;uniqueIndex;not null" json:"access_token"`
	TokenType   string         `gorm:"type:varchar(20);comment:令牌类型;default:Bearer" json:"token_type"`
	UserID      *uint          `gorm:"type:bigint;comment:用户ID;index" json:"user_id"` // 可为空，用于Client Credentials模式
	ClientID    string         `gorm:"type:varchar(100);comment:客户端ID;index;not null" json:"client_id"`
	Scope       string         `gorm:"type:varchar(500);comment:权限范围" json:"scope"`
	ExpiresAt   time.Time      `gorm:"type:datetime;comment:过期时间;index;not null" json:"expires_at"`
	Revoked     bool           `gorm:"type:tinyint(1);comment:是否已撤销;default:false" json:"revoked"`
}

func (OAuthAccessToken) TableName() string {
	return "oauth_access_tokens"
}

