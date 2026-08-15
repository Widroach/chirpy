package controller

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/internal/database"
)

func (a *API) GetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	if authorID != "" {
		id, parseErr := uuid.Parse(authorID)
		if parseErr != nil {
			responseWriter.RespondWithJson(w, 400, responseWriter.ErrorResponse{Error: "Invalid author_id"})
			return
		}
		chirps, err = a.db.GetChirpsByUserId(r.Context(), id)
	} else {
		chirps, err = a.db.GetAllChirps(r.Context())
	}

	if chirps == nil {
		chirps = []database.Chirp{}
	}

	if err != nil {
		log.Printf("error retrieving chirps: %v", err)
		responseWriter.RespondWithJson(w, 500, responseWriter.ErrorResponse{Error: "Failed to retrieve chirps."})
		return
	}

	responseWriter.RespondWithJson(w, 200, chirps)
}

func (a *API) GetChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpId")
	if id == "" {
		responseWriter.RespondWithJson(w, 400, responseWriter.ErrorResponse{Error: "You must provide a chirpId"})
		return
	}
	uuid, err := uuid.Parse(id)
	if err != nil {
		responseWriter.RespondWithJson(w, 400, responseWriter.ErrorResponse{Error: "You must provide a valid UUID"})
		return
	}
	chirp, err := a.db.GetChirp(r.Context(), uuid)
	if err != nil {
		responseWriter.RespondWithJson(w, 404, responseWriter.ErrorResponse{Error: "Chirp not found"})
		return
	}
	responseWriter.RespondWithJson(w, 200, chirp)
}
