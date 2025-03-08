package gameState

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// Service handles all game-related operations
type Service struct {
	gameStateRepo GameStateRepository
}

// NewService creates a new game service
func NewService(
	gameStateRepo GameStateRepository,
) *Service {
	return &Service{
		gameStateRepo: gameStateRepo,
	}
}

func (s *Service) CreateGameState(ctx context.Context, gameStateData GameState) (*GameState, error) {
	gameState, err := s.gameStateRepo.CreateGameState(ctx, database.CreateGameStateParams{
		ID:              gameStateData.ID,
		State:           gameStateData.State,
		GameID:          uuid.NullUUID{UUID: gameStateData.GameID},
		CurrentPlayerID: uuid.NullUUID{UUID: gameStateData.CurrentPlayerID},
		ForbiddenColor:  gameStateData.ForbiddenColor,
	})
	if err != nil {
		return nil, err
	}

	resultGameState := s.DBToModel(ctx, gameState)

	return &resultGameState, nil
}
