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
	GetBox(ctx context.Context, params database.GetBoxParams) (database.Box, error)
	ChangeBoxColor(ctx context.Context, params database.ChangeBoxColorParams) error
	SwapColors(ctx context.Context, gameID uuid.UUID, posFrom, posTo BoardPosition) error
}
