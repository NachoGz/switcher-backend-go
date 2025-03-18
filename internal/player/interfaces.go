package player

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, playerData Player) (*Player, error)
	DBToModel(ctx context.Context, dbPlayer database.Player) Player
	AssignRandomTurns(ctx context.Context, players []Player) (uuid.UUID, error)
	CountPlayers(ctx context.Context, gameID uuid.UUID) (int, error)
	GetPlayerByID(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID) (Player, error)
	GetPlayersInGame(ctx context.Context, gameID uuid.UUID) ([]Player, error)
}

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, params database.CreatePlayerParams) (database.Player, error)
	CountPlayers(ctx context.Context, gameID uuid.UUID) (int64, error)
	AssignTurnPlayer(ctx context.Context, params database.AssignTurnPlayerParams) error
	GetPlayerByID(ctx context.Context, params database.GetPlayerByIDParams) (database.Player, error)
	GetPlayersInGame(ctx context.Context, gameID uuid.UUID) ([]database.Player, error)
}
