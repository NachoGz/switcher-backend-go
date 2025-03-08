package gameState

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
)

type GameStateService interface {
	CreateGameState(ctx context.Context, gameStateData GameState) (*GameState, error)
	DBToModel(ctx context.Context, dbGameState database.GameState) GameState
}

type GameStateRepository interface {
	CreateGameState(ctx context.Context, params database.CreateGameStateParams) (database.GameState, error)
}
