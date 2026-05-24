package configs

import (
	"chatonetoone/internals/controllers/v1/chats"
	"chatonetoone/internals/controllers/v1/users"
	user_repos "chatonetoone/internals/repositories/users_repos"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func IntializeRoutes(mysqldb *gorm.DB) *mux.Router {
	router := mux.NewRouter()
	
	setupV1(router.PathPrefix("/api/v1").Subrouter(), mysqldb)

	return router
}

func setupV1(v1router *mux.Router, mysqldb *gorm.DB) {
	userRouter := v1router.PathPrefix("/users").Subrouter()
	
	repo, _ := user_repos.NewSqlUserRepository(mysqldb)
	userCtrl := users.NewUsersController(mysqldb, repo)
	users.RegisterRouters(userCtrl, userRouter)

	chatsCtrl, _ := chats.NewController()
	chats.RegisterRouters(chatsCtrl, v1router)
}