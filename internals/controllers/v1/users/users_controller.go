package users

import (
	"chatonetoone/internals/auth"
	user_repos "chatonetoone/internals/repositories/users_repos"
	"chatonetoone/internals/services/users"
	"encoding/json"
	"io"
	"log"
	"net/http"
)
type UsersController struct{
	DB 
}

func(uc *UsersController)Login(w *http.ResponseWriter, r http.Request){
	ctx := r.Context()
	select{
	case <-ctx.Done():
		log.Println("request timeout before reponse")
		return
	default:
		
	}
	var logindata map[string]any

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil{
		log.Println("error occured whiile reading body",err)
		// json.NewEncoder(w).Encode(map[string][any]{"abc":1})
	}

	if err := json.Unmarshal(body, &logindata); err != nil{
		log.Println(err)
	}

	// login service
	user_repos.NewSqlUserRepository(MysqlDB)
	users.NewLoginService(&auth.PasswordValidator{}, )
}

func Singup()
