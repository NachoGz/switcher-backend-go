package gameState

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/NachoGz/switcher-backend-go/internal/board"
	"github.com/NachoGz/switcher-backend-go/internal/figureCard"
	"github.com/NachoGz/switcher-backend-go/internal/movementCard"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

type Handlers struct {
	service             GameStateService
	playerService       player.PlayerService
	boardService        board.BoardService
	movementCardService movementCard.MovementCardService
	figureCardService   figureCard.FigureCardService
}

// NewHandlers creates a new handlers instance
func NewHandlers(service GameStateService) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting game...")
	type StartGameParams struct {
		GameID uuid.UUID `json:"game_id"`
	}

	var params StartGameParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Update state to playing
	h.service.UpdateGameState(context.Background(), params.GameID, PLAYING)

	players, err := h.playerService.GetPlayers(context.Background(), params.GameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching games", err)
		return
	}

	// Assign random turns
	firstPlayerID, err := h.playerService.AssignRandomTurns(context.Background(), players)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error setting turns", err)
		return
	}

	// Set current player
	h.service.UpdateCurrentPlayer(context.Background(), params.GameID, firstPlayerID)

	// Create board and decks for players
	err = h.boardService.ConfigureBoard(context.Background(), params.GameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error configuring board", err)
		return
	}

	err = h.movementCardService.CreateMovementDeck(context.Background(), params.GameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating movement cards deck", err)
		return
	}

	err = h.figureCardService.CreateFigureCardDeck(context.Background(), params.GameID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating figure cards deck", err)
		return
	}
}
