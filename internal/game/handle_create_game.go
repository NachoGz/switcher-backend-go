package game

import (
	"encoding/json"
	"log"
	"net/http"

	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/player"

	"github.com/NachoGz/switcher-backend-go/internal/utils"
)

// Handlers struct holds handlers with service dependency
type Handlers struct {
	service GameService
}

// NewHandlers creates a new handlers instance
func NewHandlers(service GameService) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating game")

	type GameRequest struct {
		Game   Game          `json:"game"`
		Player player.Player `json:"player"`
	}

	var params GameRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Use service to create game
	game, newGameState, newPlayer, err := h.service.CreateGame(r.Context(), params.Game, params.Player)
	if err != nil || game == nil || newGameState == nil || newPlayer == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating game", err)
		return
	}

	// Build response
	type Response struct {
		Game      Game                `json:"game"`
		GameState gameState.GameState `json:"game_state"`
		Player    player.Player       `json:"player"`
	}

	response := Response{
		Game:      *game,
		GameState: *newGameState,
		Player:    *newPlayer,
	}

	utils.RespondWithJSON(w, http.StatusCreated, response)
}
