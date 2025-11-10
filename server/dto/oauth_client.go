package dto

import "gorm.io/datatypes"

type CreateOAuthClientRequest struct {
	ClientSecret string         `json:"client_secret" validate:"required"`
	Name         string         `json:"name" validate:"required,min=3,max=20"`
	Description  string         `json:"description" validate:"omitempty,max=255"`
	Logo         string         `json:"logo" validate:"omitempty,url"`
	RedirectURIs datatypes.JSON `json:"redirect_uris" validate:"required"`
	GrantTypes   datatypes.JSON `json:"grant_types" validate:"required"`
	Scopes       datatypes.JSON `json:"scopes" validate:"required"`
	Status       int            `json:"status" validate:"required,oneof=1 0"`
}

type OAuthClientListResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Status int    `json:"status"`
}
