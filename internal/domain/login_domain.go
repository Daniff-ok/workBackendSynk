package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenResponce struct {
	AccessToken  string `json:"access_token"`
	Expires      uint   `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type AccountResponce struct {
	ID           string `json:"id"`
	Name         string `json:"username"`
	Avatar       string `json:"avatar"`
	Access       bool   `json:"access"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserID string
	Access bool
	jwt.RegisteredClaims
}
