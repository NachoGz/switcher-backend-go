package gameState

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/google/uuid"
)

// Service handles all game-related operations
type Service struct {
	gameStateRepo GameStateRepository
	playerRepo    player.PlayerRepository
}

// NewService creates a new game service
func NewService(
	gameStateRepo GameStateRepository,
	playerRepo player.PlayerRepository,

) *Service {
	return &Service{
		gameStateRepo: gameStateRepo,
		playerRepo:    playerRepo,
	}
}

// Ensure Service implements GameStateService
var _ GameStateService = (*Service)(nil)

func (s *Service) CreateGameState(ctx context.Context, gameStateData GameState) (*GameState, error) {
	gameState, err := s.gameStateRepo.CreateGameState(ctx, database.CreateGameStateParams{
		ID:              gameStateData.ID,
		State:           string(gameStateData.State),
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

func (s *Service) UpdateGameState(ctx context.Context, gameID uuid.UUID, state State) error {
	err := s.gameStateRepo.UpdateGameState(ctx, database.UpdateGameStateParams{
		GameID: uuid.NullUUID{UUID: gameID},
		State:  string(state),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateCurrentPlayer(ctx context.Context, gameID uuid.UUID, currentPlayerID uuid.UUID) error {
	err := s.gameStateRepo.UpdateCurrentPlayer(ctx, database.UpdateCurrentPlayerParams{
		GameID:          uuid.NullUUID{UUID: gameID},
		CurrentPlayerID: uuid.NullUUID{UUID: currentPlayerID},
	})
	if err != nil {
		return err
	}

	return nil
}
