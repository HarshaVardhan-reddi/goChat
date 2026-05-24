package chats

import (
	"chatonetoone/internals/helpers"
	"chatonetoone/internals/ws"
	"log"
	"net/http"
)

type ChatsController struct{

}

func NewController() (*ChatsController, error){
	return &ChatsController{}, nil
}

func (cc *ChatsController) InitiateChat(w http.ResponseWriter, r *http.Request){
	id := ws.Identifier("1")
	conn, err := ws.NewConnection(id, w, r)
	if err != nil{
		log.Println("Err in webscoket upgrader", err)
		helpers.RespondWithError(w,http.StatusInternalServerError, "failed to upgrade websocket")
	}
	ws.ChatHub.PutConnection(conn)
}