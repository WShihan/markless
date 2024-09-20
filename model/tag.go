package model

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	CreateTime time.Time `json:"create_time"`
	UserID     int       `json:"user_id"`
	Links      []Link    `gorm:"many2many:link_tags;"`
}
