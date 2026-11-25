package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                  string `json:"id"`
	Name                string `json:"global_name"`
	Avatar              string `json:"avatar"`
	Access              bool
	Acces_token         string
	Refresh_token       string
	Expires_in          uint
	Local_refresh_token string
}
