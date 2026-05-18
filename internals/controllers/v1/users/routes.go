package users

import "github.com/gorilla/mux"

func RegisterRouters(uc *UsersController, router *mux.Router){
	router.HandleFunc("/login", uc.Login).Methods("POST")
	router.HandleFunc("/signup", uc.Singup).Methods("POST")
}