package chats

import "github.com/gorilla/mux"

func RegisterRouters(uc *ChatsController, router *mux.Router){
	router.HandleFunc("/startchat", uc.InitiateChat)
}

