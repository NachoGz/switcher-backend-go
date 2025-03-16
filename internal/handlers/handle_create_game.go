package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/game"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/player"

	"github.com/NachoGz/switcher-backend-go/internal/utils"
)

func (h *GameHandlers) HandleCreateGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating game")

	type GameRequest struct {
		Game   game.Game     `json:"game"`
		Player player.Player `json:"player"`
	}

	var params GameRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Use service to create game
	newGame, newGameState, newPlayer, err := h.gameService.CreateGame(r.Context(), params.Game, params.Player)
	if err != nil || newGame == nil || newGameState == nil || newPlayer == nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating game", err)
		return
	}

	// Build response
	type Response struct {
		Game      game.Game           `json:"game"`
		GameState gameState.GameState `json:"game_state"`
		Player    player.Player       `json:"player"`
	}

	response := Response{
		Game:      *newGame,
		GameState: *newGameState,
		Player:    *newPlayer,
	}

	utils.RespondWithJSON(w, http.StatusCreated, response)
}
