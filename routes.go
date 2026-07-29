package main

import (
	"net/http"

	"github.com/widroach/chirpy/http/config"
	"github.com/widroach/chirpy/http/controller"
)

func RegisterRoutes(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /admin/metrics", controller.HitsController(cfg))
	mux.HandleFunc("POST /admin/reset", controller.ResetController(cfg))

	mux.HandleFunc("GET /api/healthz", controller.HealthzController)
	mux.HandleFunc("POST /api/users", controller.CreateNewUser(cfg))

	mux.HandleFunc("GET /api/chirps", controller.GetAllChirps(cfg))
	mux.HandleFunc("GET /api/chirps/{chirpId}", controller.GetChirp(cfg))
	mux.HandleFunc("POST /api/chirps", cfg.JWTMiddleware(controller.CreateNewChirp(cfg)))
	mux.HandleFunc("POST /api/login", controller.LoginController(cfg))
}
