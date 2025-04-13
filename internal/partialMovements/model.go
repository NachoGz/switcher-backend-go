package partialMovements

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type PartialMovement struct {
	ID             uuid.UUID `json:"id"`
	PosFromX       int       `json:"pos_from_x"`
	PosFromY       int       `json:"pos_from_y"`
	PosToX         int       `json:"pos_to_x"`
	PosToY         int       `json:"pos_to_y"`
	GameID         uuid.UUID `json:"game_id"`
	PlayerID       uuid.UUID `json:"player_id"`
	MovementCardID uuid.UUID `json:"movement_card_id"`
}

// DBToModel converts a database game to a model game with player count
func (s *Service) DBToModel(ctx context.Context, dbPartialMovement database.PartialMovement) PartialMovement {
	return PartialMovement{
		ID:             dbPartialMovement.ID,
		PosFromX:       int(dbPartialMovement.PosFromX),
		PosFromY:       int(dbPartialMovement.PosFromY),
		PosToX:         int(dbPartialMovement.PosToX),
		PosToY:         int(dbPartialMovement.PosToY),
		GameID:         dbPartialMovement.GameID,
		PlayerID:       dbPartialMovement.PlayerID,
		MovementCardID: dbPartialMovement.MovementCardID,
	}
}
