package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	MAX_CHIRP = 140
)

func ValidateController(w http.ResponseWriter, r *http.Request) {
	type RequestMessage struct {
		Body string `json:"body"`
	}
	type ResponseSuccess struct {
		Valid bool `json:"valid"`
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

	RespondWithJson(w, 200, ResponseSuccess{Valid: true})
}
