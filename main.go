package main

import (
	"log"
	"net/http"

	"github.com/widroach/chirpy/controller"
)

var (
	cfg = controller.ApiConfig{}
)

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	RegisterRoutes(serveMux)
	log.Fatal(server.ListenAndServe())
}
