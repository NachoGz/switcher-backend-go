package gamestate_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
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
	return args.Error(1)
}

func (m *MockGameStateRepository) UpdateCurrentPlayer(ctx context.Context, params database.UpdateCurrentPlayerParams) error {
	args := m.Called(ctx, params)
	return args.Error(1)
}
