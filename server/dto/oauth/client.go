package oauthdto

import (
	"time"

	"gorm.io/datatypes"
)

type CreateOAuthClientRequest struct {
	// 必填密钥字段
	ClientSecret       string `json:"client_secret" validate:"required"`
	AccessTokenSecret  string `json:"access_token_secret" validate:"required"`
	RefreshTokenSecret string `json:"refresh_token_secret" validate:"required"`
	SubjectSecret      string `json:"subject_secret" validate:"required"`

	// 必填基本字段
	Name         string         `json:"name" validate:"required,min=3,max=20"`
	RedirectURIs datatypes.JSON `json:"redirect_uris" validate:"required"`
	GrantTypes   datatypes.JSON `json:"grant_types" validate:"required"`
	Scopes       datatypes.JSON `json:"scopes" validate:"required"`
	Status       int            `json:"status" validate:"required,oneof=1 0"`

	// 可选基本字段
	Description string `json:"description" validate:"omitempty,max=255"`
	Logo        string `json:"logo" validate:"omitempty,url"`

	// 可选配置字段（不传则后端用默认值，单位：秒）
	AuthCodeExpire     *int    `json:"auth_code_expire" validate:"omitempty,min=60,max=600"`
	AccessTokenExpire  *int    `json:"access_token_expire" validate:"omitempty,min=300,max=86400"`
	RefreshTokenExpire *int    `json:"refresh_token_expire" validate:"omitempty,min=3600,max=31536000"`
	SubjectLength      *int    `json:"subject_length" validate:"omitempty,min=8,max=64"`
	SubjectPrefix      *string `json:"subject_prefix" validate:"omitempty,max=20"`
}

type OAuthClientListResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Status int    `json:"status"`
}

type OAuthClientDetailResponse struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Logo         string         `json:"logo"`
	RedirectURIs datatypes.JSON `json:"redirect_uris"`
	GrantTypes   datatypes.JSON `json:"grant_types"`
	Scopes       datatypes.JSON `json:"scopes"`
	Status       int            `json:"status"`

	// 配置字段（不暴露密钥，单位：秒）
	AuthCodeExpire     int    `json:"auth_code_expire"`
	AccessTokenSecret  string `json:"-"`
	AccessTokenExpire  int    `json:"access_token_expire"`
	RefreshTokenSecret string `json:"-"`
	RefreshTokenExpire int    `json:"refresh_token_expire"`
	SubjectSecret      string `json:"-"`
	SubjectLength      int    `json:"subject_length"`
	SubjectPrefix      string `json:"subject_prefix"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateOAuthClientRequest struct {
	// 必填基本字段
	Name string `json:"name" validate:"required,min=3,max=20"`

	// 可选基本字段
	Description  *string         `json:"description" validate:"omitempty,max=255"`
	Logo         *string         `json:"logo" validate:"omitempty,url"`
	RedirectURIs *datatypes.JSON `json:"redirect_uris" validate:"omitempty"`
	GrantTypes   *datatypes.JSON `json:"grant_types" validate:"omitempty"`
	Scopes       *datatypes.JSON `json:"scopes" validate:"omitempty"`
	Status       *int            `json:"status" validate:"omitempty,oneof=1 0"`

	// 可选密钥字段（用于轮换，不传则不更新）
	ClientSecret       *string `json:"client_secret" validate:"omitempty"`
	AccessTokenSecret  *string `json:"access_token_secret" validate:"omitempty"`
	RefreshTokenSecret *string `json:"refresh_token_secret" validate:"omitempty"`
	SubjectSecret      *string `json:"subject_secret" validate:"omitempty"`

	// 可选配置字段（单位：秒）
	AuthCodeExpire     *int    `json:"auth_code_expire" validate:"omitempty,min=60,max=600"`
	AccessTokenExpire  *int    `json:"access_token_expire" validate:"omitempty,min=300,max=86400"`
	RefreshTokenExpire *int    `json:"refresh_token_expire" validate:"omitempty,min=3600,max=31536000"`
	SubjectLength      *int    `json:"subject_length" validate:"omitempty,min=8,max=64"`
	SubjectPrefix      *string `json:"subject_prefix" validate:"omitempty,max=20"`
}
