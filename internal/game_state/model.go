package gameState

import (
	"context"
	"database/sql"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

// State enum
const (
	PLAYING  string = "PLAYING"
	WAITING  string = "WAITING"
	FINISHED string = "FINISHED"
)

type GameState struct {
	ID              uuid.UUID
	State           string
	GameID          uuid.UUID
	CurrentPlayerID uuid.UUID
	ForbiddenColor  sql.NullString
}

// DBToModel converts a database game state to a model game state
func (s *Service) DBToModel(ctx context.Context, dbGameState database.GameState) GameState {
	return GameState{
		ID:              dbGameState.ID,
		State:           dbGameState.State,
		GameID:          dbGameState.GameID.UUID,
		CurrentPlayerID: dbGameState.CurrentPlayerID.UUID,
		ForbiddenColor:  dbGameState.ForbiddenColor,
	}
}
