package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/game"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/middleware"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

	// Create database queries
	dbQueries := database.New(dbConn)

	// Create repositories
	gameRepo := game.NewGameRepository(dbQueries)
	gameStateRepo := gameState.NewGameStateRepository(dbQueries)
	playerRepo := player.NewPlayerRepository(dbQueries)

	// Create services
	gameStateService := gameState.NewService(gameStateRepo, playerRepo)
	playerService := player.NewService(playerRepo)
	gameService := game.NewService(gameRepo, gameStateRepo, playerRepo, gameStateService, playerService)

	// Create handlers
	gameHandlers := game.NewHandlers(gameService)

	// Configure routes
	mux := http.NewServeMux()

	// Game routes
	mux.HandleFunc("POST /games", gameHandlers.HandleCreateGame)
	mux.HandleFunc("GET /games", gameHandlers.HandleGetGames)
	mux.HandleFunc("GET /games/{gameID}", gameHandlers.HandleGetGameByID)

	// Add middleware
	handler := middleware.CORSMiddleware(mux)

	// Start server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
