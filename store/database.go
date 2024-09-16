package store

import (
	"markee/logging"
	"markee/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(databaseURL string) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		logging.Logger.Error(err.Error())
		logging.Logger.Fatal("failed to open database: %v", err)
	}
	DB = db
	db.Logger.LogMode(logger.Info)
	// 迁移 schema
	er := DB.AutoMigrate(&model.Link{}, &model.User{}, &model.Tag{}, &model.Category{})
	if err != nil {
		logging.Logger.Fatal("failed to migrate: %v", er)
	}
}
