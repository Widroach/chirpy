package controller

import (
	"log"
	"net/http"

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
