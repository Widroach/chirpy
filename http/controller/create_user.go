package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/widroach/chirpy/http/config"
	"github.com/widroach/chirpy/internal/auth"
	"github.com/widroach/chirpy/internal/database"
)

type postRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateNewUser(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := postRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("error decoding request %v", err)
			RespondWithJson(w, 400, struct{ Error string }{Error: "Request is invalid, it should contain email and password field"})
			return
		}

		hashedPassword, err := auth.HashPassword(request.Password)
		if err != nil {
			log.Printf("error creating new user %v", err)
			RespondWithJson(w, 500, struct{ Error string }{Error: "Failed to create the user"})
			return
		}

		user, err := cfg.Db.CreateUser(r.Context(), database.CreateUserParams{Email: request.Email, HashedPassword: hashedPassword})
		if err != nil {
			log.Printf("error creating new user %v", err)
			RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong while creating new user"})
			return
		}
	
		RespondWithJson(w, 201, user)
	}
}
