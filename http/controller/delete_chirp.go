package controller

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	responseWriter "github.com/widroach/chirpy/http"
	"github.com/widroach/chirpy/http/middleware"
)

func (a *API) DeleteChirpById(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.USER_ID_KEY)
	chirpId, err := uuid.Parse(r.PathValue("chirpId"))
	if err != nil {
		responseWriter.RespondWithJson(w, 404, responseWriter.ErrorResponse{Error: "Invalid chirp ID"})
		return
	}

	entry, err := a.db.GetChirp(r.Context(), chirpId)
	if err != nil {
		responseWriter.RespondWithJson(w, 404, responseWriter.ErrorResponse{Error: "Chirp not found"})
		return
	}

	if entry.UserID != userId {
		responseWriter.RespondWithJson(w, 403, responseWriter.ErrorResponse{Error: "The requested chirp does not belong to you"})
		return
	}

	_, err = a.db.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		log.Printf("error deleting the chirp from the database: %v", err)
		responseWriter.RespondWithJson(w, 404, responseWriter.ErrorResponse{Error: "Something went wrong"})
		return
	}

	w.WriteHeader(204)
}
