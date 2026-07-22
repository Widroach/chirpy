package config

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/widroach/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
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
	w.Write([]byte(text))
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) RegisterDatabase(dbQueries *database.Queries) {
	cfg.Db = dbQueries
}
