package store

import (
	"markless/model"
	"testing"
)

func TestLinkExist(t *testing.T) {
	InitDB("../markless.db")
	user := model.User{}
	DB.First(&user)
	t.Log("初始化数据库连接")
	if LinkExist("https://github.com/umputun/remark42", user) {
		t.Log("link exist")
	}
}

func TestLinkGetByUser(t *testing.T) {
	InitDB("../markless.db")
	user := model.User{}
	DB.First(&user)
	t.Log("初始化数据库连接")
	link, err := GetLinkByUser(user, 1)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(link)
}
