package chats

import (
	"chatonetoone/internals/helpers"
	"chatonetoone/internals/models"
	user_repos "chatonetoone/internals/repositories/users_repos"
	"chatonetoone/internals/services"
	"chatonetoone/internals/services/chats"
	"chatonetoone/internals/services/users"
	"chatonetoone/internals/ws"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type ChatsController struct{
	db *gorm.DB
}

func NewController(db *gorm.DB) (*ChatsController, error){
	return &ChatsController{db: db}, nil
}

func (cc *ChatsController) InitiateChat(w http.ResponseWriter, r *http.Request){
	repo, _ := user_repos.NewSqlUserRepository(cc.db)
	ctx, errInSettingContext := users.SetCurrentUserContext(r, repo)
	if(errInSettingContext != nil){
		helpers.RespondWithError(w,http.StatusInternalServerError,errInSettingContext.Error())
		return
	}

	user, ok := (*ctx).Value("currentuser").(*models.User)
	if !ok {
		helpers.RespondWithError(w, http.StatusInternalServerError, "user not found in context")
		return
	}

	conn, err := ws.NewConnection( ws.Identifier( strconv.Itoa(int(user.ID))), w, r)

	if err != nil{
		log.Println("Err in webscoket upgrader", err)
		helpers.RespondWithError(w,http.StatusInternalServerError, "failed to upgrade websocket")
		return
	}

	err = ws.ChatHub.PutConnection(conn)
	if err != nil {
		log.Println("Err in putting connection to hub", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "failed to register connection")
		return
	}

	ep := services.NewEventProcessor(conn, ws.ChatHub)

	chats.StartNewMessagingService(conn, ws.ChatHub, ep)
}
