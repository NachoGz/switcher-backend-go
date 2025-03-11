package movementCards

import "github.com/google/uuid"

type TypeEnum string

const (
	LINEAR_CONT   TypeEnum = "linear_cont"
	LINEAR_SPA    TypeEnum = "linear_esp"
	DIAGONAL_CONT TypeEnum = "diagonal_cont"
	DIAGONAL_SPA  TypeEnum = "diagonal_spa"
	L_RIGHT       TypeEnum = "l_right"
	L_LEFT        TypeEnum = "l_left"
	LINEAR_LAT    TypeEnum = "linear_lat"
)

type MovementCard struct {
	ID          uuid.UUID `json:"id"`
	Type        TypeEnum  `json:"type"`
	Description string    `json:"description"`
	Used        bool      `json:"used"`
	PlayerID    uuid.UUID `json:"player_id"`
	GameID      uuid.UUID `json:"game_id"`
}
