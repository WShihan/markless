package store

import (
	"fmt"
	"markless/model"
	"markless/util"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(databaseURL string) {
	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{})
	if err != nil {
		util.Logger.Error(err.Error())
		util.Logger.Fatal(fmt.Sprintf("failed to open database: %v", err))
	}
	DB = db
	db.Logger.LogMode(logger.Info)
	// 迁移 schema
	er := DB.AutoMigrate(&model.Link{}, &model.User{}, &model.Tag{}, &model.Archive{})
	if err != nil {
		util.Logger.Fatal(fmt.Sprintf("failed to migrate: %v", er))
	}
}
