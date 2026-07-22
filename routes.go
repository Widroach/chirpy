package main

import (
	"net/http"

	"github.com/widroach/chirpy/controller"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", controller.HealthzController)
	mux.HandleFunc("GET /admin/metrics", cfg.HitsController)
	mux.HandleFunc("POST /admin/reset", cfg.ResetHitsController)
	mux.HandleFunc("POST /api/validate_chirp", controller.ValidateController)
}
