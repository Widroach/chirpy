package controller

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func RespondWithJson(w http.ResponseWriter, statusCode int, payload any) {
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(payloadJson)
}
