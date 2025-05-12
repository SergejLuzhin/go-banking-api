package service

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT создает токен с userID в поле subject
func GenerateJWT(userID int64, secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10), // правильно сохраняем userID как строку числа
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
