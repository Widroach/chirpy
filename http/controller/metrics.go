// http/controller/metrics.go
package controller

import (
	"log"
	"net/http"

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
			RespondWithJson(w, 403, struct{ Error string }{Error: "Forbidden action"})
			return
		}
		if err := cfg.Db.DeleteAllUsers(r.Context()); err != nil {
			log.Printf("error deleting users %v", err)
			RespondWithJson(w, 500, struct{ Error string }{Error: "Failed to delete the users"})
			return
		}
		cfg.FileserverHits.Store(0)

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}
}

func MiddlewareMetricsInc(cfg *config.ApiConfig) func(http.Handler) http.Handler {
	return cfg.MiddlewareMetricsInc
}
