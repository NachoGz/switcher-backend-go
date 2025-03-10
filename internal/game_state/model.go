package gameState

import (
	"context"
	"database/sql"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// State enum
type State string

const (
	PLAYING  State = "playing"
	WAITING  State = "waiting"
	FINISHED State = "finished"
)

type GameState struct {
	ID              uuid.UUID
	State           State
	GameID          uuid.UUID
	CurrentPlayerID uuid.UUID
	ForbiddenColor  sql.NullString
}

// DBToModel converts a database game state to a model game state
func (s *Service) DBToModel(ctx context.Context, dbGameState database.GameState) GameState {
	return GameState{
		ID:              dbGameState.ID,
		State:           State(dbGameState.State),
		GameID:          dbGameState.GameID.UUID,
		CurrentPlayerID: dbGameState.CurrentPlayerID.UUID,
		ForbiddenColor:  dbGameState.ForbiddenColor,
	}
}
