package movementCards

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type MovementCardRepository interface {
	CreateMovementCard(ctx context.Context, params database.CreateMovementCardParams) (database.MovementCard, error)
	GetMovementDeck(ctx context.Context, gameID uuid.UUID) ([]database.MovementCard, error)
	AssignMovementCard(ctx context.Context, params database.AssignMovementCardParams) error
}

type MovementCardService interface {
	CreateMovementDeck(ctx context.Context, gameID uuid.UUID) error
}
