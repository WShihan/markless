package store

import (
	"markless/model"
	"testing"
)

func TestTagStat(t *testing.T) {
	InitDB("../markless.db")
	user := model.User{}
	err := DB.First(&user).Error
	staMap := TagStat(user)
	if err != nil {
		t.Errorf(err.Error())
	}
	t.Log(staMap)
}
