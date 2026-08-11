package middleware

import (
	"context"
	"log"
	"net/http"

	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/internal/auth"
)

type JWT struct {
	secret string
}

func NewJWT(secret string) *JWT {
	return &JWT{secret: secret}
}

func (j *JWT) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Missing authorization header"})
			return
		}
		userId, err := auth.ValidateJWT(token, j.secret)
		if err != nil {
			log.Printf("error validating JWT token: %v", err)
			responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Invalid JWT in the Authorization: Bearer header"})
			return
		}
		ctx := context.WithValue(r.Context(), USER_ID_KEY, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
