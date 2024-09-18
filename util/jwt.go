package util

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Middleware func(http.Handler) http.Handler

var (
	secretKey = []byte("secretKeyffff")
)

func CreateJWT(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string) (uid string, err error) {
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
		exp := claims["exp"]
		fmt.Println("uid:", uid)
		fmt.Println("Expires at:", exp)
	} else {
		return "", fmt.Errorf("invalid token")
	}

	return uid, nil
}
