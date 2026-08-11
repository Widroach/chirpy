package middleware

import (
	"context"
	"net/http"

	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/internal/auth"
)

func ValidateRefreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Missing authorization header"})
			return
		}
		ctx := context.WithValue(r.Context(), REFRESH_TOKEN, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
