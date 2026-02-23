package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    RoleUser,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (uint, string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil })
	if err != nil || token.Valid == false {
		return 0, "", errors.New("Token Not Valid")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok != true {
		return 0, "", errors.New("Claims Not Valid")
	}
	userID, ok := claims["user_id"].(float64)
	if ok != true {
		return 0, "", errors.New("UserID Not Valid")
	}
	role, ok := claims["role"].(string)
	if ok != true {
		return 0, "", errors.New("Role Not Valid")
	}
	return uint(userID), role, nil
}
