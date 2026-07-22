package controller

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/widroach/chirpy/http/config"
	"github.com/widroach/chirpy/internal/database"
)

func GetAllChirps(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := cfg.Db.GetAllChirps(r.Context())
		if err != nil {
			log.Printf("error retrieving chirps: %v", err)
			RespondWithJson(w, 500, ErrorResponse{Error: "Failed to retrieve chirps."})
			return
		}
		if chirps == nil {
			chirps = []database.Chirp{}
		}
		RespondWithJson(w, 200, chirps)
	}
}

func GetChirp(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("chirpId")
		if id == "" {
			RespondWithJson(w, 400, ErrorResponse{Error: "You must provide a chirpId"})
			return
		}
		uuid, err := uuid.Parse(id)
		if err != nil {
			RespondWithJson(w, 400, ErrorResponse{Error: "You must provide a valid UUID"})
			return
		}
		chirp, err := cfg.Db.GetChirp(r.Context(), uuid)
		if err != nil {
			RespondWithJson(w, 404, ErrorResponse{Error: "Chirp not found"})
			return
		}
		RespondWithJson(w, 200, chirp)
	}
}
