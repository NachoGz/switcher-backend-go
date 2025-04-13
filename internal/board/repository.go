package board

import (
	"context"
	"database/sql"

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

// GetBox fetches a box for a specific game in the given position
func (r *PostgresBoardRepository) GetBox(ctx context.Context, params database.GetBoxParams) (database.Box, error) {
	return r.queries.GetBox(ctx, params)
}

// SwitchBoxes switches the color of two colors
func (r *PostgresBoardRepository) ChangeBoxColor(ctx context.Context, params database.ChangeBoxColorParams) error {
	return r.queries.ChangeBoxColor(ctx, params)
}

// WithTx returns a new instance of the BoardRepository that uses transaction provided by sqlc
func (r *PostgresBoardRepository) WithTx(tx *sql.Tx) BoardRepository {
	return &PostgresBoardRepository{
		queries: r.queries.WithTx(tx),
	}
}
