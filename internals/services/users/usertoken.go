package users

import (
	"chatonetoone/internals/auth"
	user_repos "chatonetoone/internals/repositories/users_repos"

	"context"
	"errors"
	"net/http"
)

func DecodeAuthToken(r *http.Request) (int,error) {
	bt := r.Header.Get("authorization") // baerer token
	if(bt == ""){
		return -1, errors.New("auth token is missing in the request")		
	}
	jwthandler := auth.NewJwtHandler("")
	uid, err := jwthandler.Decode(bt)
	if err != nil{
		return -1 ,err
	}
	return int(uid), nil
}

func SetCurrentUserContext(r *http.Request, user_repo user_repos.UserRepository) (*context.Context, error){
	uid, err := DecodeAuthToken(r)
	if err != nil{
		return nil, err
	}
	user, errInFindingUser := user_repo.FindUserByID(uid)
	if(errInFindingUser != nil){
		return nil, errInFindingUser
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx,"currentuser",user)
	return &ctx, nil
}