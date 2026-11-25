package service

import (
	"context"
	"fmt"
	"os"
	"time"
	"crypto/rand"
	"encoding/base64"

	"rndmBackend/internal/domain"
	"rndmBackend/internal/gateway"
	"rndmBackend/internal/repo"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	accessTTL = 2 * time.Hour
)

type LoginService struct {
	gate      *gateway.Logingateway
	repo      *repo.LoginRepo
	jwtSecret []byte
}

func NewLoginService(gate *gateway.Logingateway, repo *repo.LoginRepo) *LoginService {
	return &LoginService{
		gate: gate,
		repo: repo,
		jwtSecret: []byte(os.Getenv("JWT_SECRET")),
	}
}

func (s *LoginService) Login(code string, ctx context.Context) (domain.AccountResponce, error) {
	tokens, err := s.gate.GetAccessToken(code)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	user, err := s.gate.GetAccountInfo(tokens.AccessToken)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	user.Acces_token = tokens.AccessToken
	user.Expires_in = tokens.Expires
	user.Refresh_token = tokens.RefreshToken

	// DELETE
	user.Access = true
	// DELETE

	_, err = s.repo.GetAccById(user.ID, ctx)
	if err == gorm.ErrRecordNotFound {
		err = s.repo.CreateUser(user, ctx)
		if err != nil {
			return domain.AccountResponce{}, err
		}
	} else if err != nil {
		return domain.AccountResponce{}, err
	} else {
		err := s.repo.UpdateUser(user, ctx)
		if err != nil {
			return domain.AccountResponce{}, err
		}
	}

	accessToken, err := generateAccessToken(user.ID, user.Access, s.jwtSecret)
	if err != nil {
		return domain.AccountResponce{}, fmt.Errorf("Generate jwt error: %w", err)
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return domain.AccountResponce{}, fmt.Errorf("Generate refresh token error: %w", err)
	}

	accountResponce := domain.AccountResponce{
		ID: user.ID,
		Name: user.Name,
		Avatar: user.Avatar,
		Access: user.Access,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return accountResponce, nil
}

func generateAccessToken(userId string, access bool, secret []byte) (string, error) {
	now := time.Now()

	claims := domain.Claims{
		UserID: userId,
		Access: access,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil

}