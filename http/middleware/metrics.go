package middleware

import (
	"net/http"
	"sync/atomic"
)

type Metrics struct {
	hits atomic.Int32
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) Inc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (m *Metrics) Reset() {
	m.hits.Store(0)
}

func (m *Metrics) Value() int32 {
	return m.hits.Load()
}
