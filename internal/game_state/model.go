package gameState

import (
	"database/sql"

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
