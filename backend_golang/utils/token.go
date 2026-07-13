// this file is used to generate new JWT Token
package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// getJWTSecret mengambil secret dari environment variable JWT_SECRET.
// Aplikasi menolak jalan tanpa secret agar tidak memakai nilai default yang tidak aman.
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET tidak diset. Set environment variable JWT_SECRET terlebih dahulu.")
	}
	return []byte(secret)
}

// GenerateJWT membuat token JWT untuk userID dengan masa berlaku sesuai duration.
func GenerateJWT(userID uint, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
	})

	// Sign token with secret key
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Validator
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Chek if that login with HMAC Methods?
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
