package game

import (
	"context"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/google/uuid"
)

type GameService interface {
	CreateGame(ctx context.Context, gameData Game, playerData player.Player) (*Game, *gameState.GameState, *player.Player, error)
	GetAvailableGames(ctx context.Context, numPlayers int, page int, limit int, name string) ([]Game, int, error)
	GetGameByID(ctx context.Context, id uuid.UUID) (*Game, error)
	DeleteGame(ctx context.Context, id uuid.UUID) error
}

type GameRepository interface {
	CreateGame(ctx context.Context, params database.CreateGameParams) (database.Game, error)
	GetAvailableGames(ctx context.Context) ([]database.Game, error)
	GetGameById(ctx context.Context, id uuid.UUID) (database.Game, error)
	DeleteGame(ctx context.Context, id uuid.UUID) error
}
