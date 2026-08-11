package main

import (
	"net/http"

	"github.com/widroach/chirpy/http/controller"
	"github.com/widroach/chirpy/http/middleware"
)

func RegisterRoutes(mux *http.ServeMux, api *controller.API, jwt *middleware.JWT) {
	mux.Handle("/app/", api.Metrics().Inc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /admin/metrics", api.HitsController)
	mux.HandleFunc("POST /admin/reset", api.ResetController)

	mux.HandleFunc("GET /api/healthz", controller.HealthzController)
	mux.HandleFunc("POST /api/users", api.CreateNewUser)

	mux.HandleFunc("GET /api/chirps", api.GetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", api.GetChirp)
	mux.Handle("POST /api/chirps", jwt.Authenticate(http.HandlerFunc(api.CreateNewChirp)))
	mux.HandleFunc("POST /api/login", api.LoginController)
}
