package model

import (
	"time"

	"gorm.io/gorm"
)

type Archive struct {
	gorm.Model
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	LinkID     int       `json:"link_id" gorm:"unique"`
	UpdateTime time.Time `json:"update_time"`
}
