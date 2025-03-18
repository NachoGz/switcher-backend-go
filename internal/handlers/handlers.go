package handlers

import (
	"github.com/NachoGz/switcher-backend-go/internal/board"
	"github.com/NachoGz/switcher-backend-go/internal/figureCard"
	"github.com/NachoGz/switcher-backend-go/internal/game"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/movementCard"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/NachoGz/switcher-backend-go/internal/websocket"
)

type GameStateHandlers struct {
	gameStateService    gameState.GameStateService
	playerService       player.PlayerService
	boardService        board.BoardService
	movementCardService movementCard.MovementCardService
	figureCardService   figureCard.FigureCardService
	wsHub               websocket.WebSocketHub
}

// NewHandlers creates a new handlers instance
func NewGameStateHandlers(gameStateService gameState.GameStateService, playerService player.PlayerService,
	boardService board.BoardService, movementCardService movementCard.MovementCardService,
	figureCardService figureCard.FigureCardService, wsHub websocket.WebSocketHub) *GameStateHandlers {
	return &GameStateHandlers{
		gameStateService:    gameStateService,
		playerService:       playerService,
		boardService:        boardService,
		movementCardService: movementCardService,
		figureCardService:   figureCardService,
		wsHub:               wsHub,
	}
}

// Handlers struct holds handlers with service dependency
type GameHandlers struct {
	gameService game.GameService
	wsHub       websocket.WebSocketHub
}

// NewHandlers creates a new handlers instance
func NewGameHandlers(gameService game.GameService, wsHub websocket.WebSocketHub) *GameHandlers {
	return &GameHandlers{
		gameService: gameService,
		wsHub:       wsHub,
	}
}

// Handlers struct holds handlers with seriveces dependancies
type PlayerHandlers struct {
	playerService    player.PlayerService
	gameService      game.GameService
	gameStateService gameState.GameStateService
	wsHub            websocket.WebSocketHub
}

// NewHandlers creates a new handlers instance
func NewPlayerHandlers(playerService player.PlayerService, gameService game.GameService, gameStateService gameState.GameStateService, wsHub websocket.WebSocketHub) *PlayerHandlers {
	return &PlayerHandlers{
		playerService:    playerService,
		gameService:      gameService,
		gameStateService: gameStateService,
		wsHub:            wsHub,
	}
}
