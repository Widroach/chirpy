package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var (
	cfg = apiConfig{}
)

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	serveMux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	serveMux.HandleFunc("GET /healthz", healthzController)
	serveMux.HandleFunc("GET /metrics", cfg.hitsController)
	serveMux.HandleFunc("POST /reset", cfg.resetHitsController)

	log.Fatal(server.ListenAndServe())
}

func healthzController(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) hitsController(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.Write([]byte(text))
}

func (cfg *apiConfig) resetHitsController(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
