package figureCard

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type FigureCardService interface {
	CreateFigureCardDeck(ctx context.Context, gameID uuid.UUID) error
	DBToModel(ctx context.Context, dbFigureCard database.FigureCard) FigureCard
}

type FigureCardRepository interface {
	CreateFigureCard(ctx context.Context, params database.CreateFigureCardParams) (database.FigureCard, error)
}
