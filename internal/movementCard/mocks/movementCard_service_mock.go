package movementCard_mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockMovementCardService struct {
	mock.Mock
}

func (m *MockMovementCardService) CreateMovementCardDeck(ctx context.Context, gameID uuid.UUID) error {
	args := m.Called(ctx, gameID)
	return args.Error(0)
}
