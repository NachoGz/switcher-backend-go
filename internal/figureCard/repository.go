package figureCard

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
)

// PostgresFigureCardRepository implements FigureCardRepository for Postgres
type PostgresFigureCardRepository struct {
	queries *database.Queries
}

// FigureCardRepository creatres a new FigureCard repository
func NewFigureCardRepository(queries *database.Queries) FigureCardRepository {
	return &PostgresFigureCardRepository{
		queries: queries,
	}
}

// CreateFigureCard creates a figure card
func (r *PostgresFigureCardRepository) CreateFigureCard(ctx context.Context, params database.CreateFigureCardParams) (database.FigureCard, error) {
	return r.queries.CreateFigureCard(ctx, params)
}
