package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/http/middleware"
	"github.com/widroach/chirpy/internal/auth"
	"github.com/widroach/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
	Secret         string
}

const METRICS_PAGE = `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`

func (cfg *ApiConfig) HitsController(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf(METRICS_PAGE, cfg.FileserverHits.Load())
	if _, err := w.Write([]byte(text)); err != nil {
		log.Printf("error failed to respond to client: %v", err)
	}
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) JWTMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Missing authorization header"})
			return
		}
		userId, err := auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			log.Printf("error validating JWT token: %v", err)
			responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Invalid JWT in the Authorization: Bearer header"})
			return
		}
		ctx := context.WithValue(r.Context(), middleware.UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *ApiConfig) RegisterDatabase(dbQueries *database.Queries) {
	cfg.Db = dbQueries
}
