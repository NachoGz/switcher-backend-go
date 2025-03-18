package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

func (h *GameStateHandlers) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting game...")
	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	// Update state to playing
	err = h.gameStateService.UpdateGameState(context.Background(), gameID, gameState.PLAYING)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating game state", err)
		return
	}

	players, err := h.playerService.GetPlayersInGame(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching players", err)
		return
	}

	// Assign random turns
	firstPlayerID, err := h.playerService.AssignRandomTurns(context.Background(), players)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error setting turns", err)
		return
	}

	// Set current player
	err = h.gameStateService.UpdateCurrentPlayer(context.Background(), gameID, firstPlayerID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating current player", err)
		return
	}

	// Create board and decks for players
	err = h.boardService.ConfigureBoard(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error configuring board", err)
		return
	}

	err = h.movementCardService.CreateMovementCardDeck(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating movement card deck", err)
		return
	}

	err = h.figureCardService.CreateFigureCardDeck(context.Background(), gameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating figure card deck", err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Game started successfully",
	})

	h.wsHub.BroadcastEvent(uuid.Nil, "GAMES_LIST_UPDATE")
	h.wsHub.BroadcastEvent(uuid.Nil, fmt.Sprintf("%s:GAME_STARTED", gameID))

}
