package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func WithJJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Print("failed to authenticate user")
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Print("failed to authenticate user")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userId := claims["userID"].(string)

		_, err = store.GetUserById(userId)
		if err != nil {
			log.Print("failed to authenticate user")
			permissionDenied(w)
			return
		}

		handlerFunc.ServeHTTP(w, r)
	}
}

func permissionDenied(w http.ResponseWriter) {
	WriteJson(w, http.StatusUnauthorized, ErrorResponse{Error: fmt.Errorf("permission denied").Error()})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func validateJWT(t string) (*jwt.Token, error) {
	secret := Envs.JWTSecret

	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func CreateJWT(s []byte, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(id)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	tokenString, err := token.SignedString(s)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
