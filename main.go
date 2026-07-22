package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

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
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to db %v", err)
	}
	dbQueries := database.New(db)
	cfg.RegisterDatabase(dbQueries)
	cfg.Platform = strings.ToLower(os.Getenv("PLATFORM"))

	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	RegisterRoutes(serveMux, &cfg)
	log.Fatal(server.ListenAndServe())
}
