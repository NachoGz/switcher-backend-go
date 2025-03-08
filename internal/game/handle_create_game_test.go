package game_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NachoGz/switcher-backend-go/internal/game"
	game_mock "github.com/NachoGz/switcher-backend-go/internal/game/mocks"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreateGame_Success(t *testing.T) {
	// Setup mock
	mockService := new(game_mock.MockGameService)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	playerID := uuid.New()
	password := "secret"

	requestGame := game.Game{
		Name:       "Test Game",
		MaxPlayers: 4,
		MinPlayers: 2,
		IsPrivate:  true,
		Password:   &password,
	}

	requestPlayer := player.Player{
		Name: "Test Player",
		Host: true,
	}

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

	responsePlayer := player.Player{
		ID:          playerID,
		Name:        "Test Player",
		Host:        true,
		GameID:      gameID,
		GameStateID: gameStateID,
	}

	// Setup expectations
	mockService.On("CreateGame", mock.Anything, requestGame, requestPlayer).
		Return(&responseGame, &responseGameState, &responsePlayer, nil)

	// Create handlers with mock service
	handlers := game.NewHandlers(mockService)

	// Create request body
	requestBody := map[string]interface{}{
		"game":   requestGame,
		"player": requestPlayer,
	}

	reqBodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleCreateGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "game")
	assert.Contains(t, response, "game_state")
	assert.Contains(t, response, "player")

	// Verify mock was called
	mockService.AssertExpectations(t)
}

func TestHandleCreateGame_InvalidRequestBody(t *testing.T) {
	// Setup mock
	mockService := new(game_mock.MockGameService)

	// Create handlers with mock service
	handlers := game.NewHandlers(mockService)

	// Create invalid request body
	reqBodyBytes := []byte(`{invalid json}`)
	req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleCreateGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid request body", response["error"])

	// Ensure service was not called
	mockService.AssertNotCalled(t, "CreateGame")
}

func TestHandleCreateGame_ServiceError(t *testing.T) {
	// Setup mock
	mockService := new(game_mock.MockGameService)

	// Test data
	requestGame := game.Game{
		Name:       "Test Game",
		MaxPlayers: 4,
		MinPlayers: 2,
	}

	requestPlayer := player.Player{
		Name: "Test Player",
		Host: true,
	}

	// Create empty objects to return
	emptyGame := &game.Game{}
	emptyGameState := &gameState.GameState{}
	emptyPlayer := &player.Player{}

	// Setup expectations with error, return empty objects plus an error
	mockService.On("CreateGame", mock.Anything, requestGame, requestPlayer).
		Return(emptyGame, emptyGameState, emptyPlayer, errors.New("service error"))

	// Create handlers with mock service
	handlers := game.NewHandlers(mockService)

	// Create request body
	requestBody := map[string]interface{}{
		"game":   requestGame,
		"player": requestPlayer,
	}

	reqBodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/games", bytes.NewReader(reqBodyBytes))
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleCreateGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Error creating game", response["error"])

	// Verify mock was called
	mockService.AssertExpectations(t)
}
