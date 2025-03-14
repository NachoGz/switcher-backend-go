package figureCard_mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockFigureCardRepository struct {
	mock.Mock
}

func (m *MockFigureCardRepository) CreateFigureCard(ctx context.Context, gameID uuid.UUID) error {
	args := m.Called(ctx, gameID)
	return args.Error(1)
}
