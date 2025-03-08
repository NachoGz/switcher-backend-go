package player_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPlayerRepository struct {
	mock.Mock
}

// CreatePlayer creates a new player
func (m *MockPlayerRepository) CreatePlayer(ctx context.Context, params database.CreatePlayerParams) (database.Player, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.Player), args.Error(1)
}

// CountPlayers counts players for a game
func (m *MockPlayerRepository) CountPlayers(ctx context.Context, gameID uuid.NullUUID) (int64, error) {
	args := m.Called(ctx, gameID)
	return args.Get(0).(int64), args.Error(1)
}
