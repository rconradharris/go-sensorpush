package sensorpush

import (
	"context"
	"net/http"
)

type AuthService service

//
// Authorization
//

type Authorization string

type authorizeRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authorizeResponse struct {
	Authorization string `json:authorization`
}

func (s *AuthService) Authorize(ctx context.Context, email, password string) (Authorization, error) {
	a0 := Authorization("")

	areq := authorizeRequest{
		Email:    email,
		Password: password,
	}

	req, err := s.c.NewBaseRequest(ctx, http.MethodPost, "oauth/authorize", areq)
	if err != nil {
		return a0, err
	}

	aresp := authorizeResponse{}
	_, err = s.c.Do(req, &aresp)
	if err != nil {
		return a0, err
	}

	return Authorization(aresp.Authorization), nil
}

//
// AccessToken
//

type AccessToken string

type accessTokenRequest struct {
	Authorization string `json:"authorization"`
}

type accessTokenResponse struct {
	AccessToken string `json:"accesstoken"`
}

func (s *AuthService) AccessToken(ctx context.Context, auth Authorization) (AccessToken, error) {
	a0 := AccessToken("")

	areq := accessTokenRequest{
		Authorization: string(auth),
	}

	req, err := s.c.NewBaseRequest(ctx, http.MethodPost, "oauth/accesstoken", areq)
	if err != nil {
		return a0, err
	}

	aresp := accessTokenResponse{}
	_, err = s.c.Do(req, &aresp)
	if err != nil {
		return a0, err
	}

	return AccessToken(aresp.AccessToken), nil
}
