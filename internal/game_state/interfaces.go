package gameState

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type GameStateService interface {
	CreateGameState(ctx context.Context, gameStateData GameState) (*GameState, error)
	DBToModel(ctx context.Context, dbGameState database.GameState) GameState
	UpdateGameState(ctx context.Context, gameID uuid.UUID, state State) error
	UpdateCurrentPlayer(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) error
}

type GameStateRepository interface {
	CreateGameState(ctx context.Context, params database.CreateGameStateParams) (database.GameState, error)
	UpdateGameState(ctx context.Context, params database.UpdateGameStateParams) error
	UpdateCurrentPlayer(ctx context.Context, params database.UpdateCurrentPlayerParams) error
}
