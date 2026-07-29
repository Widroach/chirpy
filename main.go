package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/widroach/chirpy/http/config"
	"github.com/widroach/chirpy/internal/database"
)

var (
	cfg = config.ApiConfig{
		FileserverHits: atomic.Int32{},
	}
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error reading .env: %v", err)
	}

	jwtSecret := os.Getenv("SECRET")
	cfg.Secret = strings.ToLower(jwtSecret)

	platform := os.Getenv("PLATFORM")
	cfg.Platform = strings.ToLower(platform)

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to db %v", err)
	}
	dbQueries := database.New(db)
	cfg.RegisterDatabase(dbQueries)

	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler:           serveMux,
		Addr:              ":8080",
		ReadHeaderTimeout: 2 * time.Second,
	}

	RegisterRoutes(serveMux, &cfg)
	log.Fatal(server.ListenAndServe())
}
