package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/widroach/chirpy/http/config"
)

type postRequest struct {
	Email string `json:"email"`
}

func CreateNewUser(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := postRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			RespondWithJson(w, 400, struct{ Error string }{Error: "Request should contain email field"})
			return
		}
		user, err := cfg.Db.CreateUser(r.Context(), request.Email)
		if err != nil {
			log.Printf("error creating new user %v", err)
			RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong while creating new user"})
			return
		}
		userJson, err := json.Marshal(user)
		if err != nil {
			RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong"})
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(201)
		w.Write(userJson)
	}
}
