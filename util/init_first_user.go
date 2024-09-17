package util

import (
	"markee/model"
	"markee/store"
)

func InitAdmin(username string, password string) {
	user := model.User{}
	store.DB.Find(&user, "username = ?", username)
	if user.Username == "" || user.Password == "" {
		user.Username = username
		user.Password = password
		user.Admin = true
		store.DB.Create(&user)
	}
}
