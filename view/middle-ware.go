package view

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

type Middleware func(http.Handler) http.Handler

var (
	secretKey = []byte("secretKeyffff")
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s in %v", r.URL.Path, time.Since(start))
	})
}

func CreateJWT(msg string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "msg",
		"exp":      jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Username:", claims["username"])
		fmt.Println("Expires at:", claims["exp"])
	} else {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func Protect(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		var jwt = ""
		if authHeader == "" {
			tokenCookie, trr := r.Cookie("markee-token")
			if trr != nil && authHeader == "" {
				Redirect(w, r, "/login")
				return
			}
			jwt = tokenCookie.Value

		} else {
			jwt = authHeader[len("Bearer "):]
		}

		if jwt == "" {
			Redirect(w, r, "/login")
			return
		}

		_, err := validateJWT(jwt)
		if err != nil {
			Redirect(w, r, "/login")
			return
		}
		next(w, r, ps)
	}
}

func ApplyMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
