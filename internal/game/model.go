package game

import "github.com/google/uuid"

type Game struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	MaxPlayers int       `json:"max_players"`
	MinPlayers int       `json:"min_players"`
	IsPrivate  bool      `json:"is_private"`
	Password   *string   `json:"password"`
}
