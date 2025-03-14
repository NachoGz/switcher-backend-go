package game

import (
	"context"
	"database/sql"
	"strings"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/NachoGz/switcher-backend-go/internal/utils"
	"github.com/google/uuid"
)

// Service handles all game-related operations
type Service struct {
	gameRepo         GameRepository
	gameStateRepo    gameState.GameStateRepository
	playerRepo       player.PlayerRepository
	gameStateService gameState.GameStateService
	playerService    player.PlayerService
}

// NewService creates a new game service
func NewService(
	gameRepo GameRepository,
	gameStateRepo gameState.GameStateRepository,
	playerRepo player.PlayerRepository,
	gameStateService gameState.GameStateService,
	playerService player.PlayerService,
) *Service {
	return &Service{
		gameRepo:         gameRepo,
		gameStateRepo:    gameStateRepo,
		playerRepo:       playerRepo,
		gameStateService: gameStateService,
		playerService:    playerService,
	}
}

// Ensure Service implements GameService
var _ GameService = (*Service)(nil)

// GetAvailableGames gets all available games with player counts
func (s *Service) GetAvailableGames(ctx context.Context, numPlayers int, page int, limit int, name string) ([]Game, int, error) {
	dbGames, err := s.gameRepo.GetAvailableGames(ctx)
	if err != nil {
		return nil, 0, err
	}

	var filteredGames []Game

	// Convert and filter games
	for _, dbGame := range dbGames {
		game := s.DBToModel(ctx, dbGame)

		// Filter by number of players if needed
		if game.PlayersCount == numPlayers && strings.Contains(game.Name, name) {
			filteredGames = append(filteredGames, game)
		}
	}

	// Calculate total count for pagination
	totalCount := len(filteredGames)

	// Apply pagination
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	// Check bounds
	if startIndex >= len(filteredGames) {
		return []Game{}, totalCount, nil
	}

	if endIndex > len(filteredGames) {
		endIndex = len(filteredGames)
	}

	// Return paginated results
	return filteredGames[startIndex:endIndex], totalCount, nil
}

// CreateGame creates a new game with game state and first player (creator)
func (s *Service) CreateGame(ctx context.Context, gameData Game, playerData player.Player) (*Game, *gameState.GameState, *player.Player, error) {
	// Hash password if provided
	var passwordSQL sql.NullString
	if gameData.Password != nil {
		hashedPassword, err := utils.HashPassword(*gameData.Password)
		if err != nil {
			return nil, nil, nil, err
		}
		passwordSQL = sql.NullString{String: hashedPassword, Valid: true}
	} else {
		passwordSQL = sql.NullString{Valid: false}
	}

	// Create game using repository
	game, err := s.gameRepo.CreateGame(ctx, database.CreateGameParams{
		ID:         uuid.New(),
		Name:       gameData.Name,
		MaxPlayers: int32(gameData.MaxPlayers),
		MinPlayers: int32(gameData.MinPlayers),
		IsPrivate:  gameData.IsPrivate,
		Password:   passwordSQL,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// Create game state using repository
	gameStateDb, err := s.gameStateRepo.CreateGameState(ctx, database.CreateGameStateParams{
		ID:              uuid.New(),
		State:           string(gameState.WAITING),
		GameID:          game.ID,
		CurrentPlayerID: uuid.NullUUID{Valid: false},
		ForbiddenColor:  sql.NullString{Valid: false},
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// Create player using repository
	playerDb, err := s.playerRepo.CreatePlayer(ctx, database.CreatePlayerParams{
		ID:          uuid.New(),
		Name:        playerData.Name,
		Turn:        sql.NullString{String: string(playerData.Turn), Valid: playerData.Turn != ""},
		GameID:      game.ID,
		GameStateID: gameStateDb.ID,
		Host:        playerData.Host,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert to response models
	resultGame := s.DBToModel(ctx, game)

	resultGameState := s.gameStateService.DBToModel(ctx, gameStateDb)

	resultPlayer := s.playerService.DBToModel(ctx, playerDb)

	return &resultGame, &resultGameState, &resultPlayer, nil
}

// GetGameByID gets a game by its ID
func (s *Service) GetGameByID(ctx context.Context, id uuid.UUID) (*Game, error) {
	dbGame, err := s.gameRepo.GetGameById(ctx, id)
	if err != nil {
		return nil, err
	}

	game := s.DBToModel(ctx, dbGame)

	return &game, nil
}
