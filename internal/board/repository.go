package board

import (
	"context"
	"database/sql"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// PostgresBoardRepository implements BoardRepository for Postgres
type PostgresBoardRepository struct {
	db      *sql.DB
	queries *database.Queries
}

// NewGameRepository creates a new game repository
func NewBoardRepository(queries *database.Queries, db *sql.DB) BoardRepository {
	return &PostgresBoardRepository{
		db:      db,
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

// ChangeBoxColor changes the color of a box
func (r *PostgresBoardRepository) ChangeBoxColor(ctx context.Context, params database.ChangeBoxColorParams) error {
	return r.queries.ChangeBoxColor(ctx, params)
}

// SwapColors swaps the colors between two boxes
func (r *PostgresBoardRepository) SwapColors(ctx context.Context, gameID uuid.UUID, posFrom, posTo BoardPosition) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	boxFrom, err := r.queries.GetBox(ctx, database.GetBoxParams{
		GameID: gameID,
		PosX:   int32(posFrom.PosX),
		PosY:   int32(posFrom.PosY),
	})
	if err != nil {
		return err
	}

	boxTo, err := r.queries.GetBox(ctx, database.GetBoxParams{
		GameID: gameID,
		PosX:   int32(posTo.PosX),
		PosY:   int32(posTo.PosY),
	})
	if err != nil {
		return err
	}

	// Switch colors
	if err := r.queries.WithTx(tx).ChangeBoxColor(ctx, database.ChangeBoxColorParams{
		ID:    boxFrom.ID,
		Color: boxTo.Color,
	}); err != nil {
		return err
	}

	if err := r.queries.WithTx(tx).ChangeBoxColor(ctx, database.ChangeBoxColorParams{
		ID:    boxTo.ID,
		Color: boxFrom.Color,
	}); err != nil {
		return err
	}

	return tx.Commit()
}
