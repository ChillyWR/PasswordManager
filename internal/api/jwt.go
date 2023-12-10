package api

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	ExpirationTime = time.Minute * 30
)

var SigningKey = []byte("DApAJQgpjRDHa9Ad")

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"authorized": true,
		// User ID
		"user_id": userID,
		"exp":  time.Now().Add(ExpirationTime).Unix(),
	})

	tokenStr, err := token.SignedString(SigningKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
