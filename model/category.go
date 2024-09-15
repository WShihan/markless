package model

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID         int       `json:"id"`
	Name       string    `gorm:"unique" json:"name"`
	CreateTime time.Time `json:"create_time"`
	UserID     int       `json:"user_id"`
}
