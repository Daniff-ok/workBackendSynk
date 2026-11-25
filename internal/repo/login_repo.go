package repo

import (
	"context"
	"fmt"
	"rndmBackend/internal/domain"

	"gorm.io/gorm"
)

type LoginRepo struct {
	db *gorm.DB
}

func NewLoginRepo(db *gorm.DB) *LoginRepo {
	return &LoginRepo{
		db: db,
	}
}

func (r *LoginRepo) GetUserById(id string, ctx context.Context) (domain.User, error) {
	user, err := gorm.G[domain.User](r.db).Where("ID = ?", id).First(ctx)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (r *LoginRepo) GetUserByRefresh(tokenHash string, ctx context.Context) (domain.User, error) {
	user, err := gorm.G[domain.User](r.db).Where("Local_refresh_token = ?", tokenHash).First(ctx)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}


func (r *LoginRepo) CreateUser(user domain.User, ctx context.Context) error {
	err := gorm.G[domain.User](r.db).Create(ctx, &user)
	if err != nil {
		return fmt.Errorf("Create row error: %w", err)
	}
	return nil
}

func (r *LoginRepo) UpdateUser(user domain.User, ctx context.Context) error {
	_, err := gorm.G[domain.User](r.db).Where("ID = ?", user.ID).Updates(ctx, domain.User{
		Name:          user.Name,
		Avatar:        user.Avatar,
		Acces_token:   user.Acces_token,
		Refresh_token: user.Refresh_token,
	})

	if err != nil {
		return fmt.Errorf("Update row error: %w", err)
	}
	return nil
}

func (r *LoginRepo) UpdateRefreshToken(userid string, tokenHash string, ctx context.Context) error {
	_, err := gorm.G[domain.User](r.db).Where("ID = ?", userid).Update(ctx, "Local_refresh_token", tokenHash)

	if err != nil {
		return fmt.Errorf("Update refresh token error: %w", err)
	}
	return nil
}