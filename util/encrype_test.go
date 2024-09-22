package util

import (
	"markless/tool"
	"testing"
)

func TestEncrypt(t *testing.T) {
	secretKey := []byte(tool.ShortUID(32))
	uid := "1234567890"
	t.Logf("message:%s", uid)
	encrypted, err := EncryptMessage(uid, secretKey)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Logf("encrypted message:%s", encrypted)

	decryped, err := DecryptMessage(encrypted, secretKey)
	if err != nil {
		t.Logf(err.Error())
	}
	if decryped != uid {
		t.Logf("decryped message:%s", decryped)
	}
	t.Logf("decryped message:%s", decryped)
}
