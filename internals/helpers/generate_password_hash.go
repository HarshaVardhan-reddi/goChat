package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	if(len(password) <= 0 ){
		return "", errors.New("length of the password should be greater than zero")
	}
	passwordByte := []byte(password)
	rawPassHash, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if(err != nil){
		return "", err
	}
	return string(rawPassHash), nil
}