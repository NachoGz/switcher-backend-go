package game

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
)

type GameRepository interface {
	CreateGame(ctx context.Context, params database.CreateGameParams) (database.Game, error)
	GetAvailableGames(ctx context.Context) ([]database.Game, error)
}

// PostgresGameRepository implements GameRepository for Postgres
type PostgresGameRepository struct {
	queries *database.Queries
}

// NewGameRepository creates a new game repository
func NewGameRepository(queries *database.Queries) GameRepository {
	return &PostgresGameRepository{
		queries: queries,
	}
}

// CreateGame creates a new game
func (r *PostgresGameRepository) CreateGame(ctx context.Context, params database.CreateGameParams) (database.Game, error) {
	return r.queries.CreateGame(ctx, params)
}

// GetAvailableGames gets all available games
func (r *PostgresGameRepository) GetAvailableGames(ctx context.Context) ([]database.Game, error) {
	return r.queries.GetAvailableGames(ctx)
}
