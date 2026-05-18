package models

import (
	"chatonetoone/internals/helpers"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type User struct{
	gorm.Model
	Name string `json:"name"`
	PasswordHash string `json:"-"`
	Email string `json:"email"`
	Password string `json:"password" gorm:"-"`
	gorm.DeletedAt `gorm:"-"`
}

func NewUser(details json.RawMessage) (*User, error) {
	user := User{}
	if err := json.Unmarshal(details, &user); err != nil{
		return nil, err
	}
	// input validation
	if err := user.validateInput(); err != nil{
		return nil,err
	}

	hash, err := helpers.GeneratePasswordHash(user.Password)
	if err != nil{
		return nil, err
	}

	user.PasswordHash = hash
	user.Password = ""
	return &user, nil
}

func (u User) validateInput() error{
	if u.Email == ""{
		return errors.New("email filed is missing")
	}
	if u.Password == ""{
		return  errors.New("password field is missing")
	}
	if u.Name == ""{
		return errors.New("name field is missing")
	}
	return nil
}