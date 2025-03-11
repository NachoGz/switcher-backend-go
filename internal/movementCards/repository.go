package movementCards

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// PostgresMovementCardRepository implements MovementCardRepository for Postgres
type PostgresMovementCardRepository struct {
	queries *database.Queries
}

// NewMovementCardsRepository creatres a new MovementCards repository
func NewMovementCardRepository(queries *database.Queries) MovementCardRepository {
	return &PostgresMovementCardRepository{
		queries: queries,
	}
}

// CreateMovementCard creates a new movement card
func (r *PostgresMovementCardRepository) CreateMovementCard(ctx context.Context, params database.CreateMovementCardParams) (database.MovementCard, error) {
	return r.queries.CreateMovementCard(ctx, params)
}

// GetMovementDeck fetches the movement cards for a given game
func (r *PostgresMovementCardRepository) GetMovementDeck(ctx context.Context, gameID uuid.UUID) ([]database.MovementCard, error) {
	return r.queries.GetMovementDeck(ctx, gameID)
}

// AssignMovementCard assigns the movement card to the given player
func (r *PostgresMovementCardRepository) AssignMovementCard(ctx context.Context, params database.AssignMovementCardParams) error {
	return r.queries.AssignMovementCard(ctx, params)
}
