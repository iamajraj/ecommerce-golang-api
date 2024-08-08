package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
	TokenType string `json:"token_type"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generates a JWT token for a given user ID
func GenerateJWT(userID uint, tokenType string) (string, error) {
	var expirationTime time.Time
	if tokenType == "access" {
		expirationTime = time.Now().Add(15 * time.Minute)
	} else if tokenType == "refresh" {
		expirationTime = time.Now().Add(7 * 24 * time.Hour) // 7 days
	}
	claims := &Claims{
		UserID:    userID,
		TokenType: tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// validates a JWT token and returns the claims if valid
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
