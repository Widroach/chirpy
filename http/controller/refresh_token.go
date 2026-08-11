package controller

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/http/middleware"
	"github.com/widroach/chirpy/internal/auth"
	"github.com/widroach/chirpy/internal/database"
)

type Response struct {
	Token string `json:"token"`
}

func (a *API) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.REFRESH_TOKEN).(string)
	entry, err := a.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Invalid token"})
		return
	}
	if time.Now().Compare(entry.ExpiresAt.Time) >= 0 || entry.RevokedAt.Valid {
		responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Invalid token"})
		return
	}

	JWT, err := auth.MakeJWT(entry.UserID, a.secret, time.Hour)
	if err != nil {
		log.Printf("error creating the JWT token: %v", err)
		responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Something went wrong"})
		return
	}
	responseWriter.RespondWithJson(w, 200, Response{Token: JWT})
}

func (a *API) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(middleware.REFRESH_TOKEN).(string)
	_, err := a.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		responseWriter.RespondWithJson(w, 401, struct{ Error string }{Error: "Invalid token"})
		return
	}
	w.WriteHeader(204)
}

func (a *API) SaveRefreshToken(userId uuid.UUID, context context.Context) (token string, err error) {
	refreshToken := auth.MakeRefreshToken()
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    userId,
		ExpiresAt: sql.NullTime{Time: time.Now().Add(time.Hour * 24 * 60), Valid: true},
		RevokedAt: sql.NullTime{}}

	_, err = a.db.CreateRefreshToken(context, refreshTokenParams)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}
