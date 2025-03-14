package gameState

import (
	"context"
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
	gameStateService    GameStateService
	playerService       player.PlayerService
	boardService        board.BoardService
	movementCardService movementCard.MovementCardService
	figureCardService   figureCard.FigureCardService
}

// NewHandlers creates a new handlers instance
func NewHandlers(gameStateService GameStateService, playerService player.PlayerService,
	boardService board.BoardService, movementCardService movementCard.MovementCardService,
	figureCardService figureCard.FigureCardService) *Handlers {
	return &Handlers{
		gameStateService:    gameStateService,
		playerService:       playerService,
		boardService:        boardService,
		movementCardService: movementCardService,
		figureCardService:   figureCardService,
	}
}

func (h *Handlers) HandleStartGame(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting game...")
	gameID, err := uuid.Parse(r.PathValue("gameID"))
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't parse game ID", err)
		return
	}

	// Update state to playing
	err = h.gameStateService.UpdateGameState(context.Background(), gameID, PLAYING)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating game state", err)
		return
	}

	players, err := h.playerService.GetPlayers(context.Background(), gameID)
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
}
