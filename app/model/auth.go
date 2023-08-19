package model

import (
	"crypto/sha256"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type IdentityData struct {
	Username string `json:"username" validate:"required"`
	Password int    `json:"password" validate:"required"`
}

type Claims struct {
	jwt.StandardClaims
	UserID string
}

func SHA256(password, salt string) string {
	sum := sha256.Sum256([]byte(password + salt))

	return fmt.Sprintf("%x", sum)
}

type Token struct {
	Token string `json:"token"`
}
