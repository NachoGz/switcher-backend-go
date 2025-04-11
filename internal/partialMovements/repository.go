package partialMovements

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// PostgresPartialMovementRepository implements PartialMovementRepository for Postgres
type PostgresPartialMovementRepository struct {
	queries *database.Queries
}

// NewPartialMovementRepository creates a new partial movement repository
func NewPartialMovementRepository(queries *database.Queries) PartialMovementRepository {
	return &PostgresPartialMovementRepository{
		queries: queries,
	}
}

// CreatePartialMovement creates a new partial movement
func (r *PostgresPartialMovementRepository) CreatePartialMovement(ctx context.Context, params database.CreatePartialMovementParams) (database.PartialMovement, error) {
	return r.queries.CreatePartialMovement(ctx, params)
}

// UndoMovement deletes the last entry in the partial_movements table
func (r *PostgresPartialMovementRepository) UndoMovement(ctx context.Context, params database.UndoMovementParams) error {
	return r.queries.UndoMovement(ctx, params)
}

// GetPartialMovementsByPlayer fetches all the partial movements for a given player
func (r *PostgresPartialMovementRepository) GetPartialMovementsByPlayer(ctx context.Context, params database.GetPartialMovementsByPlayerParams) ([]database.PartialMovement, error) {
	return r.queries.GetPartialMovementsByPlayer(ctx, params)
}

// UndoMovementByID deletes the partial movement for the given id
func (r *PostgresPartialMovementRepository) UndoMovementByID(ctx context.Context, partialMovID uuid.UUID) error {
	return r.queries.UndoMovementByID(ctx, partialMovID)
}

// DeleteAllPartialMovementsByPlayer deletes all the partial movements for a given player
func (r *PostgresPartialMovementRepository) DeleteAllPartialMovementsByPlayer(ctx context.Context, playerID uuid.UUID) error {
	return r.queries.DeleteAllPartialMovementsByPlayer(ctx, playerID)
}
