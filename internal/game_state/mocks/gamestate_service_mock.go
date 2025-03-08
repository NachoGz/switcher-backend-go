package gamestate_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/stretchr/testify/mock"
)

type MockGameStateService struct {
	mock.Mock
}

func (m *MockGameStateService) DBToModel(ctx context.Context, dbGameState database.GameState) gameState.GameState {
	args := m.Called(ctx, dbGameState)
	return args.Get(0).(gameState.GameState)
}

func (m *MockGameStateService) CreateGameState(ctx context.Context, gameStateData gameState.GameState) (*gameState.GameState, error) {
	args := m.Called(ctx, gameStateData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gameState.GameState), args.Error(1)
}
