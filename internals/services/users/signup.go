package users

import (
	"chatonetoone/internals/auth"
	"chatonetoone/internals/models"
	user_repos "chatonetoone/internals/repositories/users_repos"
	"strconv"
	// "chatonetoone/internals/repositories/users"
)

type SignupService struct{
	repository user_repos.UserRepository
	JwtHandler *auth.JwtHandler
}

type SignupData struct{
	User models.User `json:"user"`
	AuthToken string `json:"auth_token"`
}

func NewSingupService(repo user_repos.UserRepository, jh auth.JwtHandler) *SignupService {
	return &SignupService{repository: repo, JwtHandler: &jh}
}

func(ss *SignupService) Execute(details []byte) (SignupData, error) {
	user, err := models.NewUser(details)
	sd := SignupData{}
	if err != nil{
		return sd, err
	}
	if err := ss.repository.CreateUser(user); err != nil{
		return sd, err
	}
	sd.User = *user
	sd.AuthToken, err = ss.JwtHandler.Encode(map[string]string{"email":user.Email, "id":strconv.Itoa(int(user.ID)), "name":user.FirstName}, 86400)
	if(err != nil){
		return sd, err
	}

	return sd, nil
}