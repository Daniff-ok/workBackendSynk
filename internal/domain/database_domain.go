package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                  string `json:"id"`
	Name                string `json:"global_name"`
	Avatar              string `json:"avatar"`
	Access              bool
	Acces_token         string
	Access_updated_at   time.Time
	Expires_in          uint
	Refresh_token       string
	Local_refresh_token string
}
