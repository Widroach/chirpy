package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

var (
	MAX_CHIRP = 140
	BAD_WORDS = []string{"kerfuffle", "sharbert", "fornax"}
)

func ValidateController(w http.ResponseWriter, r *http.Request) {
	type RequestMessage struct {
		Body string `json:"body"`
	}
	type CleanedMessage struct {
		CleanedBody string `json:"cleaned_body"`
	}
	type ResponseError struct {
		Error string `json:"error"`
	}

	message := RequestMessage{}
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		log.Printf("error decoding the body %v", err)
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	if len(message.Body) > MAX_CHIRP {
		log.Printf("Chirp is too long")
		RespondWithJson(w, 400, ResponseError{Error: "Chirp is too long"})
		return
	}

	filteredBody := filterMessage(message.Body)
	RespondWithJson(w, 200, CleanedMessage{CleanedBody: filteredBody})
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
