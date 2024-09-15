package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password"`
	Token     *string    `json:"token"`
	LastLogin *time.Time `json:"last_login"`
}
