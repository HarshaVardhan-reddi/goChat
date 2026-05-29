package user_repos

import "chatonetoone/internals/models"

type UserRepository interface{
	FindUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	FindUserByID(id int) (*models.User, error)
}