package configs

import (
	"chatonetoone/internals/controllers/v1/chats"
	"chatonetoone/internals/controllers/v1/users"
	user_repos "chatonetoone/internals/repositories/users_repos"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func IntializeRoutes(mysqldb *gorm.DB) http.Handler {
	router := mux.NewRouter()
	
	setupV1(router.PathPrefix("/api/v1").Subrouter(), mysqldb)

	return corsMiddleware(router)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func setupV1(v1router *mux.Router, mysqldb *gorm.DB) {
	userRouter := v1router.PathPrefix("/users").Subrouter()
	
	repo, _ := user_repos.NewSqlUserRepository(mysqldb)
	userCtrl := users.NewUsersController(mysqldb, repo)
	users.RegisterRouters(userCtrl, userRouter)

	wsrouter := v1router.PathPrefix("/ws").Subrouter()
	chatsCtrl, _ := chats.NewController(mysqldb)
	chats.RegisterRouters(chatsCtrl, wsrouter)
}