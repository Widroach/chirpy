package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/http/config"
	"github.com/widroach/chirpy/http/middleware"
	"github.com/widroach/chirpy/internal/database"
)

var (
	MAX_CHIRP = 140
	BAD_WORDS = []string{"kerfuffle", "sharbert", "fornax"}
)

type ResponseError struct {
	Error string `json:"error"`
}

type messageReq struct {
	Body string `json:"body"`
}

func CreateNewChirp(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg := messageReq{}
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			log.Printf("error decoding the request: %v", err)
			responseWriter.RespondWithJson(w, 400, struct{ Error string }{Error: "POST body doesn't contain 'body' or 'user_id'"})
			return
		}

		if len(msg.Body) > MAX_CHIRP {
			responseWriter.RespondWithJson(w, 400, ResponseError{Error: "Chirp is too long"})
			return
		}

		userId, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
		if !ok {
			responseWriter.RespondWithJson(w, 404, ResponseError{Error: "Invalid user ID"})
			return
		}
		_, err := cfg.Db.GetUserById(r.Context(), userId)
		if err != nil {
			responseWriter.RespondWithJson(w, 404, ResponseError{Error: "User not found"})
			return
		}

		filteredBody := filterMessage(msg.Body)
		chirp, err := cfg.Db.CreateChirp(r.Context(), database.CreateChirpParams{Body: filteredBody, UserID: userId})
		if err != nil {
			log.Printf("error creating chirp: %v", err)
			responseWriter.RespondWithJson(w, 500, ResponseError{Error: "Failed to create chirp"})
			return
		}
		responseWriter.RespondWithJson(w, 201, chirp)
	}
}

func filterMessage(message string) string {
	words := strings.Split(message, " ")
	filtered := []string{}

	for _, word := range words {
		wordLowerCase := strings.ToLower(word)
		lastChar := string(word[len(word)-1])
		if isBad := slices.Contains(BAD_WORDS, wordLowerCase); isBad && lastChar != "!" {
			filtered = append(filtered, strings.Repeat("*", 4))
		} else {
			filtered = append(filtered, word)
		}
	}

	return strings.Join(filtered, " ")
}
