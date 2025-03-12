package game

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type Game struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	MaxPlayers   int       `json:"max_players"`
	MinPlayers   int       `json:"min_players"`
	PlayersCount int       `json:"players_count"`
	IsPrivate    bool      `json:"is_private"`
	Password     *string   `json:"password"`
}

// DBToModel converts a database game to a model game with player count
func (s *Service) DBToModel(ctx context.Context, dbGame database.Game) Game {
	playersCount := 0
	if count, err := s.playerRepo.CountPlayers(ctx, dbGame.ID); err == nil {
		playersCount = int(count)
	}

	return Game{
		ID:           dbGame.ID,
		Name:         dbGame.Name,
		MaxPlayers:   int(dbGame.MaxPlayers),
		MinPlayers:   int(dbGame.MinPlayers),
		PlayersCount: playersCount,
		IsPrivate:    dbGame.IsPrivate,
	}
}
