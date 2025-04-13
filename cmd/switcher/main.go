package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/NachoGz/switcher-backend-go/internal/board"
	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/figureCard"
	"github.com/NachoGz/switcher-backend-go/internal/game"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/handlers"
	"github.com/NachoGz/switcher-backend-go/internal/middleware"
	"github.com/NachoGz/switcher-backend-go/internal/movementCard"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/NachoGz/switcher-backend-go/internal/websocket"
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
	boardRepo := board.NewBoardRepository(dbQueries, dbConn)
	movementCardRepo := movementCard.NewMovementCardRepository(dbQueries)
	figureCardRepo := figureCard.NewFigureCardRepository(dbQueries)

	// Create services
	gameStateService := gameState.NewService(gameStateRepo, playerRepo)
	playerService := player.NewService(playerRepo)
	gameService := game.NewService(gameRepo, gameStateRepo, playerRepo, gameStateService, playerService)
	boardService := board.NewService(boardRepo)
	movementCardService := movementCard.NewService(movementCardRepo, playerRepo)
	figureCardService := figureCard.NewService(figureCardRepo, playerRepo)

	// Create WebSocket server
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Create handlers
	gameHandlers := handlers.NewGameHandlers(gameService, playerService, wsHub)
	gameStateHandlers := handlers.NewGameStateHandlers(gameStateService, playerService, boardService, movementCardService, figureCardService, wsHub)
	playerHandlers := handlers.NewPlayerHandlers(playerService, gameService, gameStateService, wsHub)
	wsHandlers := handlers.NewWSHandlers(wsHub, gameService, playerService)

	// Configure routes
	mux := http.NewServeMux()

	// Game routes
	mux.HandleFunc("POST /games", gameHandlers.HandleCreateGame)
	mux.HandleFunc("GET /games", gameHandlers.HandleGetGames)
	mux.HandleFunc("GET /games/{gameID}", gameHandlers.HandleGetGameByID)
	mux.HandleFunc("DELETE /games/{gameID}", gameHandlers.HandleDeleteGame)
	mux.HandleFunc("GET /games/{gameID}/winner", gameHandlers.HandlerGetWinner)

	// Game State routes
	mux.HandleFunc("PATCH /game_state/start/{gameID}", gameStateHandlers.HandleStartGame)

	// Player routes
	mux.HandleFunc("POST /players/join/{gameID}", playerHandlers.HandleJoinGame)
	mux.HandleFunc("GET /players/{gameID}", playerHandlers.HandleGetPlayers)
	mux.HandleFunc("GET /players/{gameID}/{playerID}", playerHandlers.HandleGetPlayer)

	// Websocket route
	mux.HandleFunc("/ws", wsHandlers.HandleWebSocket)

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
