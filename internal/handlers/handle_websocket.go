package handlers

import (
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/game"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/NachoGz/switcher-backend-go/internal/websocket"
	"github.com/google/uuid"
)

// WSHandlers holds WebSocket handlers
type WSHandlers struct {
	hub           websocket.WebSocketHub
	gameService   game.GameService
	playerService player.PlayerService
}

// NewWSHandlers creates a new WebSocket handlers instance
func NewWSHandlers(hub websocket.WebSocketHub, gameService game.GameService, playerService player.PlayerService) *WSHandlers {
	return &WSHandlers{
		hub:           hub,
		gameService:   gameService,
		playerService: playerService,
	}
}

// HandleWebSocket handles WebSocket connections for a game
func (h *WSHandlers) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil && gameID != uuid.Nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid game ID", err)
		return
	}

	playerID, err := uuid.Parse(r.URL.Query().Get("player_id"))
	if err != nil && playerID != uuid.Nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid player ID", err)
		return
	}

	// Get the player to verify they belong to this game
	if playerID != uuid.Nil {
		player, err := h.playerService.GetPlayerByID(r.Context(), playerID, gameID)
		if err != nil || player.GameID != gameID {
			utils.RespondWithError(w, http.StatusForbidden, "Player not in this game", err)
			return
		}
	}

	log.Printf("WebSocket connection: game=%s player=%s", gameID, playerID)

	// Upgrade HTTP connection to WebSocket
	conn, err := websocket.NewConnection(w, r)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}

	// Create client
	client := &websocket.Client{
		Server:   h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		GameID:   gameID,
		PlayerID: playerID,
	}

	// Register client
	h.hub.RegisterClient(client)

	// Start client message pumps
	go client.WritePump()
	go client.ReadPump()
}
