package player

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// PostgresPlayerRepository implements PlayerRepository for Postgres
type PostgresPlayerRepository struct {
	queries *database.Queries
}

// NewPlayerRepository creates a new player repository
func NewPlayerRepository(queries *database.Queries) PlayerRepository {
	return &PostgresPlayerRepository{
		queries: queries,
	}
}

// CreatePlayer creates a new player
func (r *PostgresPlayerRepository) CreatePlayer(ctx context.Context, params database.CreatePlayerParams) (database.Player, error) {
	return r.queries.CreatePlayer(ctx, params)
}

// CountPlayers counts players for a game
func (r *PostgresPlayerRepository) CountPlayers(ctx context.Context, gameID uuid.NullUUID) (int64, error) {
	return r.queries.CountPlayers(ctx, gameID)
}

// GetPlayers fetches all the players in a game
func (r *PostgresPlayerRepository) GetPlayers(ctx context.Context, gameID uuid.NullUUID) ([]database.Player, error) {
	return r.queries.GetPlayers(ctx, gameID)
}

// AssignTurnPlayer sets the turn for the given player
func (r *PostgresPlayerRepository) AssignTurnPlayer(ctx context.Context, params database.AssignTurnPlayerParams) error {
	return r.queries.AssignTurnPlayer(ctx, params)
}
