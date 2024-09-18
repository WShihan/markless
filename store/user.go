package store

import (
	"markee/model"
)

func GetUserByUID(uid string) (user model.User, err error) {
	user = model.User{}
	err = DB.Find(&user, "uid = ?", uid).Error
	return user, err
}
