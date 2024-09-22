package tool

import (
	"markless/model"
	"testing"
	"time"
)

func TestTimeFMT(t *testing.T) {
	t.Logf("%v", TimeFMT(time.Now()))
}

func TestGetBaseTemplate(t *testing.T) {
	t.Logf("%v", GetBaseTemplate())
}

func TestExcutePath(t *testing.T) {
	t.Logf("%v", ExcutePath())
}

func TestJoinTagNames(t *testing.T) {
	tags := []model.Tag{
		{ID: 1, Name: "tag1", CreateTime: time.Now()},
		{ID: 2, Name: "tag2", CreateTime: time.Now()},
		{ID: 3, Name: "tag3", CreateTime: time.Now()},
	}
	s := JoinTagNames(tags)
	if s != "tag1&tag2&tag3" {
		panic("error")
	}
}

func TestShort_UID(t *testing.T) {
	t.Logf("short_uid: %s", ShortUID(10))
}

func TestHashMessage(t *testing.T) {
	pass, err := HashMessage("123456ff")
	if err != nil {
		t.Errorf(err.Error())
	}
	err = ValidateHash(pass, "123456ff")
	if err != nil {
		t.Errorf(err.Error())
	} else {
		t.Log("验证通过")
	}
	t.Logf("%s", pass)
}
