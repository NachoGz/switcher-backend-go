package movementCard

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type MovementCardRepository interface {
	CreateMovementCard(ctx context.Context, params database.CreateMovementCardParams) (database.MovementCard, error)
	GetMovementCardDeck(ctx context.Context, gameID uuid.UUID) ([]database.MovementCard, error)
	AssignMovementCard(ctx context.Context, params database.AssignMovementCardParams) error
}

type MovementCardService interface {
	CreateMovementCardDeck(ctx context.Context, gameID uuid.UUID) error
}
