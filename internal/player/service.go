package player

import (
	"context"
	"database/sql"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// Service handles all game-related operations
type Service struct {
	playerRepo PlayerRepository
}

// NewService creates a new game service
func NewService(
	playerRepo PlayerRepository,
) *Service {
	return &Service{
		playerRepo: playerRepo,
	}
}

func (s *Service) CreatePlayer(ctx context.Context, playerData Player) (*Player, error) {
	player, err := s.playerRepo.CreatePlayer(ctx, database.CreatePlayerParams{
		ID:          playerData.ID,
		Name:        playerData.Name,
		Turn:        sql.NullString{String: playerData.Turn},
		GameID:      uuid.NullUUID{UUID: playerData.GameID},
		GameStateID: uuid.NullUUID{UUID: playerData.GameStateID},
		Host:        playerData.Host,
	})
	if err != nil {
		return nil, err
	}

	resultPlayer := s.DBToModel(ctx, player)

	return &resultPlayer, nil
}
