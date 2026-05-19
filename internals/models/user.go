package models

import (
	"chatonetoone/internals/helpers"
	"encoding/json"
	"errors"
	"time"
)

type User struct{
	ID int64 `gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	PasswordDigest string `json:"-"`
	Email string `json:"email"`
	Password string `json:"password" gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Nickname string `json:"nickname"`
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

	user.PasswordDigest = hash
	user.Password = ""
	return &user, nil
}

func (u User) validateInput() error {
	if u.Email == ""{
		return errors.New("email filed is missing")
	}
	if u.Password == ""{
		return  errors.New("password field is missing")
	}
	// if u.FirstName == ""{
	// 	return errors.New("name field is missing")
	// }
	return nil
}