package oauthservices

import (
	"context"

	"goauth/dto/oauth"
	"goauth/services"
)

type OAuthUserInfoService struct {
	userService *services.UserService
}

func NewOAuthUserInfoService(userService *services.UserService) *OAuthUserInfoService {
	return &OAuthUserInfoService{userService: userService}
}

func (s *OAuthUserInfoService) GetUserInfo(ctx context.Context, username string) *oauthdto.UserInfoResponse {
	user, err := s.userService.GetUser(ctx, map[string]any{"username": username})
	if err != nil {
		return nil
	}

	return &oauthdto.UserInfoResponse{
		Sub:       user.Subject,
		Nickname:  user.Nickname,
		Picture:   user.Avatar,
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}
