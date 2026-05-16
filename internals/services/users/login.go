package users

import (
	"chatonetoone/internals/auth"
	"chatonetoone/internals/models"
	"chatonetoone/internals/repositories/users"
	
	"strconv"
)

// 5 days in seconds
const DEFAULT_EXPIRY_FOR_LOGIN = 259200

type LoginReponse struct{
	models.User
	AuthToken string
}

type LoginService struct{
	credentialValidator auth.CredentialValidator
	userRepository users.UserRepository
	jwtHandler *auth.JwtHandler
}

func NewLoginService(validator auth.CredentialValidator, repo users.UserRepository) *LoginService {
	return &LoginService{
		credentialValidator: validator, 
		userRepository: repo,
		jwtHandler: auth.NewJwtHandler(""),
	}
}

func(ls *LoginService) Execute(email string, password string) (LoginReponse, error) {
	user, err := ls.userRepository.FindUserByEmail(email)
	if err != nil{
		return LoginReponse{},err
	}
	// validate password
	if err := ls.credentialValidator.Validate(password, user.PasswordHash); err != nil{
		return LoginReponse{}, err
	}

	// generate jwt
	token, err := ls.jwtHandler.Encode(map[string]string{"email":user.Email, "id": strconv.Itoa(int(user.ID)), "name":user.Name}, DEFAULT_EXPIRY_FOR_LOGIN)


	return LoginReponse{User: *user, AuthToken: token}, nil
}
