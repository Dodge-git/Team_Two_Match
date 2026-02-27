package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenStr string,secret string) (uint, string, error) {

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) { 
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); ok != true{
			return nil,errors.New("isn,t correct signing method")
		}
		return []byte(secret), nil })
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
