package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/game"
	"github.com/NachoGz/switcher-backend-go/internal/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error opening database connection: %v", err)
	}
	defer dbConn.Close()

	dbQueries := database.New(dbConn)
	cfg := &middleware.ApiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	// Handlers
	mux.HandleFunc("POST /games", game.HandleCreateGame)

	// Wrap mux with CORS middleware
	handler := middleware.WithConfig(cfg)(mux)
	handler = middleware.CORSMiddleware(handler)

	srv := &http.Server{
		Addr:    ":" + port, // Listen on port 8080
		Handler: handler,
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
