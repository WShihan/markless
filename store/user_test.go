package store

import (
	"markless/model"
	"testing"
)

func TestUerGetByID(t *testing.T) {
	InitDB("../markless.db")
	user := model.User{}
	err := DB.First(&user).Error
	if err != nil {
		t.Errorf(err.Error())
	}
	_, err = GetUserByUID(user.Uid)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(user)
}
