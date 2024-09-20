package util

import (
	"testing"
)

func TestJWT(t *testing.T) {
	token, err := CreateJWT("w123456")
	if err != nil {
		panic(err)

	}
	if token == "" {
		panic("jwt 生成失败")
	}
	t.Log(token)
	uid, err := ValidateJWT(token)
	t.Log("uid:", uid)
	if err != nil {
		panic(err)
	}
	if uid != "w123456" {
		panic("jwt 验证失败")
	}
}
