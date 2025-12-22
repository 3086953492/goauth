package oauthmodels

import (
	"time"

	"gorm.io/gorm"
)

type OAuthAuthorizationCode struct {
	ID          uint           `gorm:"type:bigint;comment:授权码ID;primaryKey" json:"id"`
	CreatedAt   time.Time      `gorm:"type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime;comment:更新时间" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"type:datetime;comment:删除时间;index" json:"-"`
	Code        string         `gorm:"type:varchar(100);comment:授权码;uniqueIndex;not null" json:"code"`
	UserID      uint           `gorm:"type:bigint;comment:用户ID;index;not null" json:"user_id"`
	ClientID    string         `gorm:"type:varchar(100);comment:客户端ID;index;not null" json:"client_id"`
	RedirectURI string         `gorm:"type:varchar(500);comment:回调地址;not null" json:"redirect_uri"`
	Scope       string         `gorm:"type:varchar(500);comment:权限范围" json:"scope"`
	ExpiresAt   time.Time      `gorm:"type:datetime;comment:过期时间;index;not null" json:"expires_at"`
	Used        bool           `gorm:"type:tinyint(1);comment:是否已使用;default:false" json:"used"`
}

func (OAuthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}

