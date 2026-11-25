package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"rndmBackend/internal/domain"
	"rndmBackend/internal/gateway"
	"rndmBackend/internal/repo"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	accessTTL = 5 * time.Hour
)

type LoginService struct {
	gate      *gateway.Logingateway
	repo      *repo.LoginRepo
	jwtSecret []byte
}

func NewLoginService(gate *gateway.Logingateway, repo *repo.LoginRepo) *LoginService {
	return &LoginService{
		gate:      gate,
		repo:      repo,
		jwtSecret: []byte(os.Getenv("JWT_SECRET")),
	}
}

func (s *LoginService) Login(code string, ctx context.Context) (domain.AccountResponce, error) {
	tokens, err := s.gate.GetTokens(code)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	user, err := s.gate.GetAccountInfo(tokens.AccessToken)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	user.Acces_token = tokens.AccessToken
	user.Access_updated_at = time.Now()
	user.Expires_in = tokens.Expires
	user.Refresh_token = tokens.RefreshToken

	// DELETE
	user.Access = true
	// DELETE

	_, err = s.repo.GetUserById(user.ID, ctx)
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
	sum := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(sum[:])
	err = s.repo.UpdateUserAttribute(user.ID, "Local_refresh_token", tokenHash, ctx)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	accountResponce := domain.AccountResponce{
		ID:           user.ID,
		Name:         user.Name,
		Avatar:       user.Avatar,
		Access:       user.Access,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return accountResponce, nil
}

func (s *LoginService) Refresh(token string, ctx context.Context) (domain.AccountResponce, error) {
	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])
	userDB, err := s.repo.GetUserByRefresh(tokenHash, ctx)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	if time.Since(userDB.Access_updated_at) > time.Duration(userDB.Expires_in)*time.Second {
		newtoken, err := s.gate.RefreshAccessToken(userDB.Refresh_token)
		if err != nil {
			return domain.AccountResponce{}, err
		}
		err = s.repo.UpdateUser(domain.User{ID: userDB.ID, Acces_token: newtoken, Access_updated_at: time.Now()}, ctx)
	}

	UserApi, err := s.gate.GetAccountInfo(userDB.Acces_token)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	err = s.repo.UpdateUser(UserApi, ctx)
	if err != nil {
		return domain.AccountResponce{}, err
	}

	accessToken, err := generateAccessToken(userDB.ID, userDB.Access, s.jwtSecret)
	if err != nil {
		return domain.AccountResponce{}, fmt.Errorf("Generate jwt error: %w", err)
	}

	accountResponce := domain.AccountResponce{
		ID:          userDB.ID,
		Name:        UserApi.Name,
		Avatar:      UserApi.Avatar,
		Access:      userDB.Access,
		AccessToken: accessToken,
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
