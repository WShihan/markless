package util

import (
	"markee/model"
	"markee/store"
	"markee/tool"
)

func InitAdmin(username string, password string) {
	user := model.User{}
	store.DB.Where("username = ?", username).Find(&user)
	if user.Username == "" || user.Password == "" {
		user.Username = username
		user.Password = password
		user.Uid = tool.Short_UID(10)
		user.Admin = true
		store.DB.Create(&user)
	}
}
