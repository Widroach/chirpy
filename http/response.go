package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Printf("error failed to respond to client: %v", err)
	}
}

func RespondWithJson(w http.ResponseWriter, statusCode int, payload any) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(payloadJson); err != nil {
		log.Printf("error failed to respond to client: %v", err)
	}
}
