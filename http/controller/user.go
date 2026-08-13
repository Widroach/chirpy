package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/http/middleware"
	"github.com/widroach/chirpy/internal/auth"
	"github.com/widroach/chirpy/internal/database"
)

type RequestUpdateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *API) UpdateUser(w http.ResponseWriter, r *http.Request) {
	body := RequestUpdateUser{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		responseWriter.RespondWithJson(w, 400, struct{ Error string }{Error: "POST body doesn't contain 'email' and 'password'"})
		return
	}

	if body.Email == "" || body.Password == "" {
		responseWriter.RespondWithJson(w, 400, struct{ Error string }{Error: "POST body doesn't contain 'email' and 'password'"})
		return
	}

	userId := r.Context().Value(middleware.USER_ID_KEY).(uuid.UUID)
	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		log.Printf("error while hashing the password: %v", err)
		responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong'"})
		return

	}
	params := database.UpdateUserParams{
		ID:             userId,
		Email:          body.Email,
		HashedPassword: hashedPassword,
	}

	res, err := a.db.UpdateUser(r.Context(), params)
	if err != nil {
		responseWriter.RespondWithJson(w, 404, struct{ Error string }{Error: "User not found"})
		return
	}

	responseWriter.RespondWithJson(w, 200, res)
}
