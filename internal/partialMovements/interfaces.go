package partialMovements

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type PartialMovementService interface {
	RevertPartialMovements(ctx context.Context, gameID, playerID uuid.UUID) error
}

type PartialMovementRepository interface {
	CreatePartialMovement(ctx context.Context, params database.CreatePartialMovementParams) (database.PartialMovement, error)
	UndoMovement(ctx context.Context, params database.UndoMovementParams) error
	GetPartialMovementsByPlayer(ctx context.Context, params database.GetPartialMovementsByPlayerParams) ([]database.PartialMovement, error)
	UndoMovementByID(ctx context.Context, partialMovID uuid.UUID) error
	DeleteAllPartialMovementsByPlayer(ctx context.Context, playerID uuid.UUID) error
}
