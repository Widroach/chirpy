package controller

import (
	"fmt"
	"log"
	"net/http"

	responseWriter "github.com/widroach/chirpy/http"
)

const METRICS_PAGE = `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`

func (a *API) HitsController(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf(METRICS_PAGE, a.metrics.Value())
	if _, err := w.Write([]byte(text)); err != nil {
		log.Printf("error failed to respond to client: %v", err)
	}
}

func (a *API) ResetController(w http.ResponseWriter, r *http.Request) {
	if a.platform != "dev" {
		responseWriter.RespondWithJson(w, 403, struct{ Error string }{Error: "Forbidden action"})
		return
	}
	if err := a.db.DeleteAllUsers(r.Context()); err != nil {
		log.Printf("error deleting users %v", err)
		responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Failed to delete the users"})
		return
	}
	if err := a.db.DeleteAllChirps(r.Context()); err != nil {
		log.Printf("error deleting chirps %v", err)
		responseWriter.RespondWithJson(w, 500, struct{ Error string }{Error: "Failed to delete chirps"})
		return
	}
	a.metrics.Reset()

	responseWriter.RespondWithJson(w, http.StatusOK, struct {
		Success bool `json:"success"`
	}{Success: true})
}
