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
	if user.Username == "" || user.Password == "" {
		user.Token = &toke
		user.Username = username
		user.Password = password
		user.Uid = tool.ShortUID(10)
		user.Admin = true
		DB.Create(&user)
	}
}
