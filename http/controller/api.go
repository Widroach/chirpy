package controller

import (
	"github.com/widroach/chirpy/http/middleware"
	"github.com/widroach/chirpy/internal/database"
)

type API struct {
	db       *database.Queries
	secret   string
	platform string
	metrics  *middleware.Metrics
}

func NewAPI(db *database.Queries, secret, platform string, metrics *middleware.Metrics) *API {
	return &API{
		db:       db,
		secret:   secret,
		platform: platform,
		metrics:  metrics,
	}
}

func (a *API) Metrics() *middleware.Metrics {
	return a.metrics
}
