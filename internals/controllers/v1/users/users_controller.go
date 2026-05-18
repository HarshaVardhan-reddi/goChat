package users

import (
	"chatonetoone/internals/auth"
	"chatonetoone/internals/helpers"
	user_repos "chatonetoone/internals/repositories/users_repos"
	"chatonetoone/internals/services/users"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"gorm.io/gorm"
)

type UsersController struct {
	DB   *gorm.DB
	repo user_repos.UserRepository
}

func NewUsersController(db *gorm.DB, repo user_repos.UserRepository) *UsersController {
	return &UsersController{
		DB:   db,
		repo: repo,
	}
}

func (uc *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done():
		log.Println("request timeout before response")
		return
	default:
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error occurred while reading body", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var logindata map[string]any
	if err := json.Unmarshal(body, &logindata); err != nil {
		log.Println("error unmarshaling login data:", err)
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// login service
	// Note: You might want to assign the result of these to variables and use them
	user_repos.NewSqlUserRepository(uc.DB)
	users.NewLoginService(&auth.PasswordValidator{}, uc.repo)

	// Placeholder for successful login response
	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "login logic pending"})
}

func (uc *UsersController) Singup(w http.ResponseWriter, r *http.Request) {
	rawBody, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("error occurred while reading the request body:", err)
		helpers.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	singupservice := users.NewSingupService(uc.repo)
	user, err := singupservice.Execute(rawBody)
	if err != nil {
		log.Println("error occurred while creating the user:", err)
		helpers.RespondWithError(w, http.StatusPreconditionFailed, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, map[string]any{
		"message": "user created successfully",
		"user":    user,
	})
}
