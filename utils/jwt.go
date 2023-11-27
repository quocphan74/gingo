package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Define the secret key for signing the JWT
var jwtKey = []byte("hin")

// generateToken creates and signs a JWT token
func GenerateToken(user_id string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user_id,
		ExpiresAt: expirationTime.Unix(),
	})
	return token.SignedString(jwtKey)
}

func Parsejwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	clai := token.Claims.(*jwt.StandardClaims)
	return clai.Issuer, nil
}
