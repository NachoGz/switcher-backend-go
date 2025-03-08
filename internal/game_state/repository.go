package gameState

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
)

// PostgresGameStateRepository implements GameStateRepository for Postgres
type PostgresGameStateRepository struct {
	queries *database.Queries
}

// NewGameStateRepository creates a new game state repository
func NewGameStateRepository(queries *database.Queries) GameStateRepository {
	return &PostgresGameStateRepository{
		queries: queries,
	}
}

// CreateGameState creates a new game state
func (r *PostgresGameStateRepository) CreateGameState(ctx context.Context, params database.CreateGameStateParams) (database.GameState, error) {
	return r.queries.CreateGameState(ctx, params)
}
