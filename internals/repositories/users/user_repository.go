package users

import "chatonetoone/internals/models"

type UserRepository interface{
	FindUserByEmail(email string) (*models.User, error)
}