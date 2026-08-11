package controller

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/internal/database"
)

func (a *API) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("error retrieving chirps: %v", err)
		responseWriter.RespondWithJson(w, 500, responseWriter.ErrorResponse{Error: "Failed to retrieve chirps."})
		return
	}
	if chirps == nil {
		chirps = []database.Chirp{}
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
