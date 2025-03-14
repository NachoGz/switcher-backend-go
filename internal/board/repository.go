package board

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// PostgresBoardRepository implements BoardRepository for Postgres
type PostgresBoardRepository struct {
	queries *database.Queries
}

// NewGameRepository creates a new game repository
func NewBoardRepository(queries *database.Queries) BoardRepository {
	return &PostgresBoardRepository{
		queries: queries,
	}
}

// CreateBoard creates a new board
func (r *PostgresBoardRepository) CreateBoard(ctx context.Context, params database.CreateBoardParams) (database.Board, error) {
	return r.queries.CreateBoard(ctx, params)
}

// GetBoard fetches the board for the given ID
func (r *PostgresBoardRepository) GetBoard(ctx context.Context, gameID uuid.UUID) (database.Board, error) {
	return r.queries.GetBoard(ctx, gameID)
}

// AddBoxToBoard creates a new box within the board
func (r *PostgresBoardRepository) AddBoxToBoard(ctx context.Context, params database.AddBoxToBoardParams) (database.Box, error) {
	return r.queries.AddBoxToBoard(ctx, params)
}
