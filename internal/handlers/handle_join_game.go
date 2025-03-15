package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

func (h *PlayerHandlers) HandleJoinGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Joining game...")

	type PlayerJoinRequest struct {
		PlayerName string  `json:"player_name"`
		Password   *string `json:"password"`
	}

	var params PlayerJoinRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	game, err := h.gameService.GetGameByID(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get game with ID: %s", gameID), err)
		return
	}

	playersInGame, err := h.playerService.CountPlayers(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get amount of players in game with ID: %s", gameID), err)
	}

	if game.MaxPlayers == int(playersInGame) {
		utils.RespondWithError(w, http.StatusForbidden, "The game is full", err)
	}

	if game.IsPrivate && game.Password != nil {
		storedPasswordHash := game.Password
		// No password entered
		if params.Password == nil {
			utils.RespondWithError(w, http.StatusForbidden, "Password required for private games", err)
			return
		}

		if err := utils.CheckPasswordHash(*storedPasswordHash, *params.Password); err != nil {
			utils.RespondWithError(w, http.StatusForbidden, "Incorrect password", err)
			return
		}
	}

	gameState, err := h.gameStateService.GetGameStateByGameID(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get game state for game %s", gameID), err)
		return
	}

	// Create player
	_, err = h.playerService.CreatePlayer(context.Background(), player.Player{
		Name:        params.PlayerName,
		GameID:      gameID,
		GameStateID: gameState.ID,
		Host:        false,
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create player", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Joined game successfully",
	})
}
