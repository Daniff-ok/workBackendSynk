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
	user, err := gorm.G[domain.User](r.db).Where("id = ?", id).First(ctx)
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
	_, err := gorm.G[domain.User](r.db).Where("id = ?", user.ID).Updates(ctx, user)

	if err != nil {
		return fmt.Errorf("Update row error: %w", err)
	}
	return nil
}

func (r *LoginRepo) UpdateUserAttribute(userid string, attribute string, value any, ctx context.Context) error {
	_, err := gorm.G[domain.User](r.db).Where("id = ?", userid).Update(ctx, attribute, value)

	if err != nil {
		return fmt.Errorf("Update refresh token error: %w", err)
	}
	return nil
}