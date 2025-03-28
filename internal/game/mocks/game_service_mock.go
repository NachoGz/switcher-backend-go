package game_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/game"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockGameService struct {
	mock.Mock
}

func (m *MockGameService) CreateGame(ctx context.Context, gameData game.Game, playerData player.Player) (*game.Game, *gameState.GameState, *player.Player, error) {
	args := m.Called(ctx, gameData, playerData)
	return args.Get(0).(*game.Game), args.Get(1).(*gameState.GameState), args.Get(2).(*player.Player), args.Error(3)
}

func (m *MockGameService) GetAvailableGames(ctx context.Context, numPlayers int, page int, limit int, name string) ([]game.Game, int, error) {
	args := m.Called(ctx, numPlayers, page, limit, name)
	return args.Get(0).([]game.Game), args.Int(1), args.Error(2)
}

func (m *MockGameService) GetGameByID(ctx context.Context, id uuid.UUID) (*game.Game, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*game.Game), args.Error(1)
}

func (m *MockGameService) DeleteGame(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
