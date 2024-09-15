package store

import (
	"marky/model"

	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(databaseURL string) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	DB = db
	// 迁移 schema
	er := DB.AutoMigrate(&model.Link{}, &model.User{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", er)
	}
}
