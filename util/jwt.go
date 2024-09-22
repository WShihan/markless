package util

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type Middleware func(http.Handler) http.Handler

func CreateJWT(claims jwt.MapClaims, hmacSecret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateAndEncryptJWT(claims jwt.MapClaims, hmacSecret []byte, secretKey []byte) (string, error) {
	tokenString, err := CreateJWT(claims, hmacSecret)
	if err != nil {
		return "", err
	}
	encrypedJWT, err := EncryptMessage(tokenString, secretKey)
	if err != nil {
		return "", err
	} else {
		// Base64 编码 JWT
		return base64.StdEncoding.EncodeToString([]byte(encrypedJWT)), nil
	}
}

func ValidateJWT(tokenString string, secretKey []byte) (uid string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid = claims["uid"].(string)
	} else {
		return "", fmt.Errorf("invalid token")
	}

	return uid, nil
}

func DecryptAndVerifyJWT(tokenString string, hmacSecret []byte, secretKey []byte) (uid string, err error) {
	// Base64 解码 JWT
	decodedJWT, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return "", err
	}
	tokenString, err = DecryptMessage(string(decodedJWT), secretKey)
	if err != nil {
		return "", err
	}
	return ValidateJWT(tokenString, hmacSecret)
}
