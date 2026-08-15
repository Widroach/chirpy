package middleware

import (
	"net/http"

	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/internal/auth"
)

func VerifyPolkaApiKey(next http.Handler, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestApiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: err.Error()})
			return
		}
		if requestApiKey != apiKey {
			responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Invalid ApiKey"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
