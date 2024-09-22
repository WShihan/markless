package store

import (
	"markless/model"
	"markless/tool"
	"markless/util"
)

func InitAdmin(username string, password string) {
	user := model.User{}
	DB.Where("username = ?", username).Find(&user)
	toke, _ := util.GenerateRandomKey(12)
	pass, err := tool.HashMessage(password)
	if err != nil {
		panic("init admin error：" + err.Error())
	}
	if user.Username == "" || user.Password == "" {
		user.Token = &toke
		user.Username = username
		user.Password = pass
		user.Uid = tool.ShortUID(10)
		user.Admin = true
		DB.Create(&user)
	}
}
