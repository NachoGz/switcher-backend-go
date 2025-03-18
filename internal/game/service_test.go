package game_test

import (
	"context"
	"testing"

	"github.com/NachoGz/switcher-backend-go/internal/database"
	"github.com/NachoGz/switcher-backend-go/internal/game"
	game_mock "github.com/NachoGz/switcher-backend-go/internal/game/mocks"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	gameState_mock "github.com/NachoGz/switcher-backend-go/internal/game_state/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	player_mock "github.com/NachoGz/switcher-backend-go/internal/player/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGame(t *testing.T) {
	// Create mock repositories
	mockGameRepo := new(game_mock.MockGameRepository)
	mockGameStateRepo := new(gameState_mock.MockGameStateRepository)
	mockPlayerRepo := new(player_mock.MockPlayerRepository)

	// Create mock services for game_state and player
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)

	// Create the service with all the mocks
	service := game.NewService(
		mockGameRepo,
		mockGameStateRepo,
		mockPlayerRepo,
		mockGameStateService,
		mockPlayerService,
	)

	// Test data
	password := ""
	testGame := game.Game{Name: "Test Game", MaxPlayers: 4, MinPlayers: 2, Password: &password}
	testPlayer := player.Player{Name: "Player1"}

	// Mock responses from database
	gameID := uuid.New()
	gameStateID := uuid.New()
	playerID := uuid.New()

	// Setup database response for game
	dbGame := database.Game{
		ID:         gameID,
		Name:       "Test Game",
		MaxPlayers: 4,
		MinPlayers: 2,
		IsPrivate:  false,
	}

	// Setup database response for game state
	dbGameState := database.GameState{
		ID:     gameStateID,
		State:  string(gameState.WAITING),
		GameID: gameID,
	}

	// Setup database response for player
	dbPlayer := database.Player{
		ID:          playerID,
		Name:        "Player1",
		GameID:      gameID,
		GameStateID: gameStateID,
	}

	// Use mock.MatchedBy to match parameters regardless of UUID
	mockGameRepo.On("CreateGame", mock.Anything, mock.MatchedBy(func(params database.CreateGameParams) bool {
		return params.Name == testGame.Name &&
			params.MaxPlayers == int32(testGame.MaxPlayers) &&
			params.MinPlayers == int32(testGame.MinPlayers)
	})).Return(dbGame, nil)

	// Setup expectations for game state
	mockGameStateRepo.On("CreateGameState", mock.Anything, mock.MatchedBy(func(params database.CreateGameStateParams) bool {
		return params.State == string(gameState.WAITING) &&
			params.GameID == gameID
	})).Return(dbGameState, nil)

	// Setup expectations for player
	mockPlayerRepo.On("CreatePlayer", mock.Anything, mock.MatchedBy(func(params database.CreatePlayerParams) bool {
		return params.Name == testPlayer.Name &&
			params.GameID == gameID &&
			params.GameStateID == gameStateID
	})).Return(dbPlayer, nil)

	// Set up player count for the DBToModel call
	mockPlayerRepo.On("CountPlayers", mock.Anything, mock.MatchedBy(func(g uuid.UUID) bool {
		return g == gameID
	})).Return(int64(1), nil)

	mockGameStateService.On("DBToModel", mock.Anything, dbGameState).Return(gameState.GameState{
		ID:     gameStateID,
		State:  gameState.WAITING,
		GameID: gameID,
	})

	mockPlayerService.On("DBToModel", mock.Anything, dbPlayer).Return(player.Player{
		ID:          playerID,
		Name:        "Player1",
		GameID:      gameID,
		GameStateID: gameStateID,
	})

	// Execute the function being tested
	createdGame, createdGameState, createdPlayer, err := service.CreateGame(context.Background(), testGame, testPlayer)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, createdGame)
	assert.Equal(t, "Test Game", createdGame.Name)
	assert.Equal(t, gameID, createdGame.ID)
	assert.Equal(t, 1, createdGame.PlayersCount)

	assert.NotNil(t, createdGameState)
	assert.Equal(t, gameState.WAITING, createdGameState.State)

	assert.NotNil(t, createdPlayer)
	assert.Equal(t, "Player1", createdPlayer.Name)

	// Verify expectations
	mockGameRepo.AssertExpectations(t)
	mockGameStateRepo.AssertExpectations(t)
	mockPlayerRepo.AssertExpectations(t)
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
}
