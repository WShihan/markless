package model

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	ID         int       `json:"id"`
	Url        string    `json:"url"`
	Icon       string    `json:"icon"`
	Title      string    `json:"title"`
	Desc       string    `json:"desc"`
	Read       bool      `json:"read"`
	Archive    *Archive  `json:"archive"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time" `
	Tags       []Tag     `gorm:"many2many:link_tags;"`
	UserID     int       `json:"user_id"`
}
