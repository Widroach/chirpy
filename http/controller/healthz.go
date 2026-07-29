package controller

import (
	"net/http"

	ResponseWriter "github.com/widroach/chirpy/http"
)

func HealthzController(w http.ResponseWriter, r *http.Request) {
	ResponseWriter.RespondWithJson(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{Status: "OK"})
}
