package game_mock

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) CreateGame(ctx context.Context, params database.CreateGameParams) (database.Game, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(database.Game), args.Error(1)
}

func (m *MockGameRepository) GetAvailableGames(ctx context.Context) ([]database.Game, error) {
	args := m.Called(ctx)
	return args.Get(0).([]database.Game), args.Error(1)
}

func (m *MockGameRepository) GetGameById(ctx context.Context, id uuid.UUID) (database.Game, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Game), args.Error(1)
}
