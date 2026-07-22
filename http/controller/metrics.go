// http/controller/metrics.go
package controller

import (
	"net/http"

	"github.com/widroach/chirpy/http/config"
)

func HitsController(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.HitsController(w, r)
	}
}

func ResetHitsController(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.ResetHitsController(w, r)
	}
}

func MiddlewareMetricsInc(cfg *config.ApiConfig) func(http.Handler) http.Handler {
	return cfg.MiddlewareMetricsInc
}
