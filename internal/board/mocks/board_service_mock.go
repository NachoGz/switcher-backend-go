package board_mock

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockBoardService struct {
	mock.Mock
}

func (m *MockBoardService) ConfigureBoard(ctx context.Context, gameID uuid.UUID) error {
	args := m.Called(ctx, gameID)
	return args.Error(0)
}
