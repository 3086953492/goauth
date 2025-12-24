package oauthmodels

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OAuthClient struct {
	ID                 uint           `gorm:"type:bigint;comment:客户端ID;primaryKey" json:"id"`
	ClientSecret       string         `gorm:"type:text;comment:客户端密钥;not null" json:"-"`
	RedirectURIs       datatypes.JSON `gorm:"type:json;comment:回调地址列表" json:"redirect_uris"`
	GrantTypes         datatypes.JSON `gorm:"type:json;comment:支持的授权类型" json:"grant_types"`
	Scopes             datatypes.JSON `gorm:"type:json;comment:允许的权限范围" json:"scopes"`
	AuthCodeExpire     int            `gorm:"type:int;comment:授权码过期时间;not null" json:"auth_code_expires_at"`
	AccessTokenSecret  string         `gorm:"type:text;comment:访问令牌密钥;not null" json:"access_token_secret"`
	AccessTokenExpire  int            `gorm:"type:int;comment:访问令牌过期时间;not null" json:"access_token_expires_at"`
	RefreshTokenSecret string         `gorm:"type:text;comment:刷新令牌密钥;not null" json:"refresh_token_secret"`
	RefreshTokenExpire int            `gorm:"type:int;comment:刷新令牌过期时间;not null" json:"refresh_token_expires_at"`
	Name               string         `gorm:"type:varchar(100);comment:应用名称;not null" json:"name"`
	Description        string         `gorm:"type:text;comment:应用描述" json:"description"`
	Logo               string         `gorm:"type:varchar(500);comment:应用Logo URL" json:"logo"`
	Status             int            `gorm:"type:tinyint;comment:状态;default:1" json:"status"` // 1:启用 0:禁用
	CreatedAt          time.Time      `gorm:"type:datetime;comment:创建时间" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:datetime;comment:更新时间" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"type:datetime;comment:删除时间;index" json:"-"`
}

func (OAuthClient) TableName() string {
	return "oauth_clients"
}
