package controller

import (
	"github.com/widroach/chirpy/http/middleware"
	"github.com/widroach/chirpy/internal/database"
)

type API struct {
	db       *database.Queries
	secret   string
	PolkaKey string
	platform string
	metrics  *middleware.Metrics
}

func NewAPI(db *database.Queries, secret, polkaKey, platform string, metrics *middleware.Metrics) *API {
	return &API{
		db:       db,
		secret:   secret,
		PolkaKey: polkaKey,
		platform: platform,
		metrics:  metrics,
	}
}

func (a *API) Metrics() *middleware.Metrics {
	return a.metrics
}
