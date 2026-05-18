package configs

import (
	"chatonetoone/internals/controllers/v1/users"

	"github.com/gorilla/mux"
)

func IntializeRoutes(userCtrl *users.UsersController) *mux.Router {
	router := mux.NewRouter()
	v1router := router.PathPrefix("/api/v1").Subrouter()
	userRouter := v1router.PathPrefix("/users").Subrouter()
	
	users.RegisterRouters(userCtrl, userRouter)
	
	return router
}