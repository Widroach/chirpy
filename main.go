package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/widroach/chirpy/http/controller"
	"github.com/widroach/chirpy/http/middleware"
	"github.com/widroach/chirpy/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error reading .env: %v", err)
	}

	secret := strings.ToLower(os.Getenv("SECRET"))
	platform := strings.ToLower(os.Getenv("PLATFORM"))
	polkaKey := strings.ToLower(os.Getenv("POLKA_KEY"))

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to db %v", err)
	}
	dbQueries := database.New(db)

	api := controller.NewAPI(dbQueries, secret, polkaKey, platform, middleware.NewMetrics())

	serveMux := http.NewServeMux()
	RegisterRoutes(serveMux, api, middleware.NewJWT(secret))

	server := &http.Server{
		Handler:           serveMux,
		Addr:              ":8080",
		ReadHeaderTimeout: 2 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
