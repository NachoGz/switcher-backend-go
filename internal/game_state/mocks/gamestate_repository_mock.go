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
