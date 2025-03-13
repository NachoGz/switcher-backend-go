package board_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockBoardRepository struct {
	mock.Mock
}

func (m *MockBoardRepository) CreateBoard(ctx context.Context, params database.CreateBoardParams) (database.Board, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.Board), args.Error(1)
}

func (m *MockBoardRepository) GetBoard(ctx context.Context, gameID uuid.UUID) (database.Board, error) {
	args := m.Called(ctx, gameID)
	return args.Get(0).(database.Board), args.Error(1)
}

func (m *MockBoardRepository) AddBoxToBoard(ctx context.Context, params database.AddBoxToBoardParams) (database.Box, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.Box), args.Error(1)
}
