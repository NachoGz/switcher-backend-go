package player

import "github.com/google/uuid"

type Player struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Turn        string    `json:"turn"`
	GameID      uuid.UUID `json:"game_id"`
	GameStateID uuid.UUID `json:"game_state_id"`
	Host        bool      `json:"host"`
}
