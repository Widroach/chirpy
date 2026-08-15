package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
)

type UpgradeRequestBody struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (a *API) UpgradeUserMembership(w http.ResponseWriter, r *http.Request) {
	requestBody := UpgradeRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		responseWriter.RespondWithJson(w, 400, struct{ Error string }{Error: "Wrong request"})
		return
	}
	if requestBody.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err := a.db.GetUserById(r.Context(), requestBody.Data.UserID)
	if err != nil {
		responseWriter.RespondWithJson(w, 404, struct{ Error string }{Error: "User not found"})
		return
	}

	if err := a.db.UpgradeUserMembership(r.Context(), requestBody.Data.UserID); err != nil {
		log.Printf("error upgrading user membership: %v", err)
		responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong"})
		return
	}

	w.WriteHeader(204)
}
