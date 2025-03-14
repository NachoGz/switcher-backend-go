package gameState_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockGameStateRepository struct {
	mock.Mock
}

func (m *MockGameStateRepository) CreateGameState(ctx context.Context, params database.CreateGameStateParams) (database.GameState, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.GameState), args.Error(1)
}

func (m *MockGameStateRepository) UpdateGameState(ctx context.Context, params database.UpdateGameStateParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockGameStateRepository) UpdateCurrentPlayer(ctx context.Context, params database.UpdateCurrentPlayerParams) error {
	args := m.Called(ctx, params)
	return args.Error(0)
}

func (m *MockGameStateRepository) GetGameStateByGameID(ctx context.Context, gameID uuid.UUID) (database.GameState, error) {
	args := m.Called(ctx, gameID)
	return args.Get(0).(database.GameState), args.Error(1)
}
