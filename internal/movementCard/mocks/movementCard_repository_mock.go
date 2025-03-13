package movementCard_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockMovementCardRepository struct {
	mock.Mock
}

func (m *MockMovementCardRepository) CreateMovementCard(ctx context.Context, params database.CreateMovementCardParams) (database.MovementCard, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.MovementCard), args.Error(1)
}

func (m *MockMovementCardRepository) GetMovementDeck(ctx context.Context, gameID uuid.UUID) ([]database.MovementCard, error) {
	args := m.Called(ctx, gameID)
	return args.Get(0).([]database.MovementCard), args.Error(1)
}

func (m *MockMovementCardRepository) AssignMovementCard(ctx context.Context, params database.AssignMovementCardParams) error {
	args := m.Called(ctx, params)
	return args.Error(1)
}
