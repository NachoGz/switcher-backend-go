package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NachoGz/switcher-backend-go/internal/game"
	game_mock "github.com/NachoGz/switcher-backend-go/internal/game/mocks"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	gameState_mock "github.com/NachoGz/switcher-backend-go/internal/game_state/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/handlers"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	player_mock "github.com/NachoGz/switcher-backend-go/internal/player/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleJoinGamePublic_Success(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	playerID := uuid.New()

	// Create expected response objects
	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
	}

	responseGameState := gameState.GameState{
		ID:     gameStateID,
		State:  gameState.WAITING,
		GameID: gameID,
	}

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	responsePlayer := player.Player{
		ID:          playerID,
		Name:        "Test Player",
		Host:        false,
		GameID:      gameID,
		GameStateID: gameStateID,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&responseGame, nil)

	mockPlayerService.On("CountPlayers", mock.Anything, gameID).
		Return(responseGame.PlayersCount, nil)

	mockGameStateService.On("GetGameStateByGameID", mock.Anything, gameID).
		Return(&responseGameState, nil)

	mockPlayerService.On("CreatePlayer", mock.Anything, requestPlayer).
		Return(&responsePlayer, nil)

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Verify response body
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Equal(t, "Joined game successfully", response["message"])

	// Verify mocks are called
	mockPlayerService.AssertExpectations(t)
	mockGameService.AssertExpectations(t)
	mockGameStateService.AssertExpectations(t)
}
