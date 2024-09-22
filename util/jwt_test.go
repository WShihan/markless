package util

import (
	"markless/tool"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var claims = jwt.MapClaims{
	"uid": "1234567890",
	"exp": time.Now().Add(time.Hour * 1).Unix(),
}

func TestJWT(t *testing.T) {

	hmacSecret := []byte(tool.ShortUID(12))
	jwt, err := CreateJWT(claims, hmacSecret)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Logf("jwt:%s", jwt)

	validateJwt, err := ValidateJWT(jwt, hmacSecret)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Logf("uid:%v", validateJwt)

}
func TestEncryptAndDecriptJWT(t *testing.T) {
	hmacSecret := []byte(tool.ShortUID(12))
	secretKey := []byte(tool.ShortUID(32))
	jwt, err := CreateAndEncryptJWT(claims, hmacSecret, secretKey)
	if err != nil {
		t.Logf(err.Error())
	}
	if jwt == "" {
		t.Logf("jwt is empty")
	}
	t.Logf("encrypedJwt:%s", jwt)

	validateJwt, err := DecryptAndVerifyJWT(jwt, hmacSecret, secretKey)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Logf("uid:%v", validateJwt)
}
