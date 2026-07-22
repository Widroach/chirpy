package main

import (
	"net/http"

	"github.com/widroach/chirpy/http/config"
	"github.com/widroach/chirpy/http/controller"
)

func RegisterRoutes(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", controller.HealthzController)
	mux.HandleFunc("GET /admin/metrics", controller.HitsController(cfg))
	mux.HandleFunc("POST /admin/reset", controller.ResetHitsController(cfg))
	mux.HandleFunc("POST /api/validate_chirp", controller.ValidateController)
}
