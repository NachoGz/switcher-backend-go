package player_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPlayerService struct {
	mock.Mock
}

func (m *MockPlayerService) DBToModel(ctx context.Context, dbPlayer database.Player) player.Player {
	args := m.Called(ctx, dbPlayer)
	return args.Get(0).(player.Player)
}

func (m *MockPlayerService) CreatePlayer(ctx context.Context, playerData player.Player) (*player.Player, error) {
	args := m.Called(ctx, playerData)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*player.Player), args.Error(1)
}

func (m *MockPlayerService) GetPlayers(ctx context.Context, gameID uuid.NullUUID) ([]player.Player, error) {
	args := m.Called(ctx, gameID)
	return args.Get(0).([]player.Player), args.Error(1)
}

func (m *MockPlayerService) AssignRandomTurns(ctx context.Context, players []player.Player) (uuid.UUID, error) {
	args := m.Called(ctx, players)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
