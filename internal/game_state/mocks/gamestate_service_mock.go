package gameState_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/google/uuid"
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

func (m *MockGameStateService) UpdateGameState(ctx context.Context, gameID uuid.UUID, state gameState.State) error {
	args := m.Called(ctx, gameID, state)
	return args.Error(0)
}

func (m *MockGameStateService) UpdateCurrentPlayer(ctx context.Context, gameID uuid.UUID, currentPlayerID uuid.UUID) error {
	args := m.Called(ctx, gameID, currentPlayerID)
	return args.Error(0)
}
