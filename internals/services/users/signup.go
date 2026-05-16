package users

import (
	"chatonetoone/internals/models"
	"chatonetoone/internals/repositories/users"
)

type SignupService struct{
	repository users.UserRepository
}

func(ss *SignupService) Execute(details []byte) (*models.User, error) {
	user, err := models.NewUser(details)
	if err != nil{
		return user, err
	}
	if err := ss.repository.CreateUser(user); err != nil{
		return nil, err
	}
	return user, nil
}