package controller

import (
	"log"
	"net/http"

	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/http/config"
)

func HitsController(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.HitsController(w, r)
	}
}

func ResetController(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.Platform != "dev" {
			responseWriter.RespondWithJson(w, 403, struct{ Error string }{Error: "Forbidden action"})
			return
		}
		if err := cfg.Db.DeleteAllUsers(r.Context()); err != nil {
			log.Printf("error deleting users %v", err)
			responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Failed to delete the users"})
			return
		}
		if err := cfg.Db.DeleteAllChirps(r.Context()); err != nil {
			log.Printf("error deleting chirps %v", err)
			responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Failed to delete chirps"})
			return
		}
		cfg.FileserverHits.Store(0)

		responseWriter.RespondWithJson(w, http.StatusOK, struct {
			Success bool `json:"success"`
		}{Success: true})
	}
}

func MiddlewareMetricsInc(cfg *config.ApiConfig) func(http.Handler) http.Handler {
	return cfg.MiddlewareMetricsInc
}
