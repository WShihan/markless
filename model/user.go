package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         int        `json:"id"`
	Uid        string     `gorm:"unique" json:"uid"`
	Username   string     `gorm:"unique" json:"username"`
	Password   string     `json:"password"`
	Token      *string    `json:"token"`
	Lang       string     `json:"lang"`
	LastLogin  *time.Time `json:"last_login"`
	Categroies []Category
	Admin      bool `json:"admin"`
	Tags       []Tag
	Links      []Link
}
