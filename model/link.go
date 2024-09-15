package model

import (
	"time"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Desc       string    `json:"desc"`
	Url        string    `json:"url"`
	Read       bool      `json:"read"`
	CreateTime time.Time `json:"create_time"`
	Tags       []Tag     `gorm:"many2many:link_tags;"`
	UserID     int       `json:"user_id"`
	User       User
}
