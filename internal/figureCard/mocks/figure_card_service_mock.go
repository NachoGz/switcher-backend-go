package figure_card_mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFigureCardService struct {
	mock.Mock
}

func (m *MockFigureCardService) CreateFigureCardDeck(ctx context.Context, gameID uuid.UUID) error {
	args := m.Called(ctx, gameID)
	return args.Error(1)
}
