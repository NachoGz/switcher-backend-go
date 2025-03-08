package player

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/google/uuid"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, playerData Player) (*Player, error)
	DBToModel(ctx context.Context, dbPlayer database.Player) Player
}

type PlayerRepository interface {
	CreatePlayer(ctx context.Context, params database.CreatePlayerParams) (database.Player, error)
	CountPlayers(ctx context.Context, gameID uuid.NullUUID) (int64, error)
}
