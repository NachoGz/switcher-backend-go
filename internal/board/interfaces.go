package board

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type BoardService interface {
	ConfigureBoard(ctx context.Context, gameID uuid.UUID) error
}

type BoardRepository interface {
	CreateBoard(ctx context.Context, params database.CreateBoardParams) (database.Board, error)
	GetBoard(ctx context.Context, gameID uuid.UUID) (database.Board, error)
	AddBoxToBoard(ctx context.Context, params database.AddBoxToBoardParams) (database.Box, error)
}
