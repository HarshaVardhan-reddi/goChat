package models

import "gorm.io/gorm"

type User struct{
	gorm.Model
	Name string
	PasswordHash string
	Email string
}