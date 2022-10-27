package services

import (
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("HELLOOsecretkey12345#!@secretkeyHELLOO")

func GenerateJWT(username string) (string, error) {

	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenStr string) (string, error) {

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return "token invalid", err
	}

	if !tkn.Valid {
		return "token invalid", err
	}

	return "", err

}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
