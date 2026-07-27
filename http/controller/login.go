package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/widroach/chirpy/http/config"
	"github.com/widroach/chirpy/internal/auth"
)

type ResponseUser struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func LoginController(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request := postRequest{}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("error decoding request %v", err)
			RespondWithJson(w, 400, struct{ Error string }{Error: "Request is invalid, it should contain email and password field"})
			return
		}
		user, err := cfg.Db.GetUser(r.Context(), request.Email)
		if err != nil {
			RespondWithJson(w, 401, struct{ Error string }{Error: "Wrong username or password, try again"})
			return
		}
		ok, err := auth.CheckPasswordHash(request.Password, user.HashedPassword)
		if err != nil {
			RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong"})
			return
		}
		if !ok {
			RespondWithJson(w, 401, struct{ Error string }{Error: "Wrong username or password, try again"})
			return
		}

		response := ResponseUser{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
		RespondWithJson(w, 200, response)
	}
}
