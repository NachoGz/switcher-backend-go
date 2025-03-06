package player

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type Player struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Turn        string    `json:"turn"`
	GameID      uuid.UUID `json:"game_id"`
	GameStateID uuid.UUID `json:"game_state_id"`
	Host        bool      `json:"host"`
}

// DBToModel converts a database player to a model player
func (s *Service) DBToModel(ctx context.Context, dbPlayer database.Player) Player {
	return Player{
		ID:          dbPlayer.ID,
		Name:        dbPlayer.Name,
		Turn:        dbPlayer.Turn.String,
		GameID:      dbPlayer.GameID.UUID,
		GameStateID: dbPlayer.GameStateID.UUID,
		Host:        dbPlayer.Host,
	}
}
