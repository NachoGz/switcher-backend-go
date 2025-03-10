package player

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type Player struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Turn        TurnEnum  `json:"turn"`
	GameID      uuid.UUID `json:"game_id"`
	GameStateID uuid.UUID `json:"game_state_id"`
	Host        bool      `json:"host"`
}

// turnEnum
type TurnEnum string

const (
	FIRST  TurnEnum = "fist"
	SECOND TurnEnum = "second"
	THIRD  TurnEnum = "third"
	FOURTH TurnEnum = "fourth"
)

// DBToModel converts a database player to a model player
func (s *Service) DBToModel(ctx context.Context, dbPlayer database.Player) Player {
	return Player{
		ID:          dbPlayer.ID,
		Name:        dbPlayer.Name,
		Turn:        TurnEnum(dbPlayer.Turn.String),
		GameID:      dbPlayer.GameID.UUID,
		GameStateID: dbPlayer.GameStateID.UUID,
		Host:        dbPlayer.Host,
	}
}
