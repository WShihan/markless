package store

import (
	"markless/model"
	"testing"
)

func TestInitDB(t *testing.T) {
	InitDB("../markless.db")
	err := DB.First(&model.User{}).Error
	if err != nil {
		panic("数据库错误：" + err.Error())
	}
}
