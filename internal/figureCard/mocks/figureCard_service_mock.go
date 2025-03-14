package figureCard_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/figureCard"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFigureCardService struct {
	mock.Mock
}

func (m *MockFigureCardService) CreateFigureCardDeck(ctx context.Context, gameID uuid.UUID) error {
	args := m.Called(ctx, gameID)
	return args.Error(0)
}

func (m *MockFigureCardService) DBToModel(ctx context.Context, dbFigureCard database.FigureCard) figureCard.FigureCard {
	args := m.Called(ctx, dbFigureCard)
	return args.Get(0).(figureCard.FigureCard)
}
