package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/internal/auth"
)

type LoginRequest struct {
	postRequest
	ExpiresInSeconds int `json:"expires_in_seconds"`
}

type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (a *API) LoginController(w http.ResponseWriter, r *http.Request) {
	request := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("error decoding request %v", err)
		responseWriter.RespondWithJson(w, 400, struct{ Error string }{Error: "Request is invalid, it should contain email and password field"})
		return
	}
	user, err := a.db.GetUser(r.Context(), request.Email)
	if err != nil {
		responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Wrong username or password, try again"})
		return
	}
	ok, err := auth.CheckPasswordHash(request.Password, user.HashedPassword)
	if err != nil {
		responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong"})
		return
	}
	if !ok {
		responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Wrong username or password, try again"})
		return
	}

	if request.ExpiresInSeconds <= 0 || request.ExpiresInSeconds > int(time.Hour.Seconds()) {
		request.ExpiresInSeconds = 3600
	}
	token, err := auth.MakeJWT(user.ID, a.secret, time.Duration(request.ExpiresInSeconds)*time.Second)
	if err != nil {
		log.Printf("error making JWT: %v", err)
		responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong"})
		return
	}

	response := LoginResponse{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Token: token}
	responseWriter.RespondWithJson(w, 200, response)
}
