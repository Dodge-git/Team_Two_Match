package auth

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password),10/*bcrypt.DefaultCost*/)
	if err != nil {
		return "",err
	}
	return string(hashed),nil
}
func CheckPassword(password string, hashedPassword string)error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(password))
}
func ValidatePassword(password string)error{
	if utf8.RuneCountInString(password) < 8 {
		return errors.New("invalid password")
	}

	for _, r := range password {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r){
			return errors.New("invalid password")
		}
	}
	return nil

}