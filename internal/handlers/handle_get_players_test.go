package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	game_mock "github.com/NachoGz/switcher-backend-go/internal/game/mocks"
	gameState_mock "github.com/NachoGz/switcher-backend-go/internal/game_state/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/handlers"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	player_mock "github.com/NachoGz/switcher-backend-go/internal/player/mocks"
	websocket_mock "github.com/NachoGz/switcher-backend-go/internal/websocket/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleGetPlayer_Success(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Create test games
	gameID := uuid.New()
	gameStateID := uuid.New()

	responsePlayers := []player.Player{
		{
			ID:          uuid.New(),
			Name:        "Test Player 1",
			Host:        true,
			GameID:      gameID,
			GameStateID: gameStateID,
		},
		{
			ID:          uuid.New(),
			Name:        "Test Player 2",
			Host:        false,
			GameID:      gameID,
			GameStateID: gameStateID,
		},
	}
	// Set up mock expectations
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return(responsePlayers, nil)

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/players/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetPlayers(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response []player.Player
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Equal(t, response, responsePlayers)

	// Verify player array
	assert.Equal(t, 2, len(responsePlayers))

	// Verift mock was called
	mockPlayerService.AssertExpectations(t)
}

func TestHandleGetPlayers_InvalidGameID(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/players/", nil)
	req.SetPathValue("gameID", "invalid-id")
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetPlayers(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't parse game ID", response["error"])

	// Verify mock was called
	mockPlayerService.AssertNotCalled(t, "GetPlayersInGame")
}

func TestHandleGetPlayers_ServiceError(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()

	// Setup expectations
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return([]player.Player{}, errors.New("error fetching players"))

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/players/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetPlayers(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "failed to fetch players", response["error"])

	// Verify mock was called
	mockPlayerService.AssertExpectations(t)
}

func TestHandleGetPlayerByID_Success(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Create test games
	gameID := uuid.New()
	playerID := uuid.New()
	gameStateID := uuid.New()

	responsePlayer := player.Player{
		ID:          playerID,
		Name:        "Test Player",
		Host:        true,
		GameID:      gameID,
		GameStateID: gameStateID,
	}
	// Set up mock expectations
	mockPlayerService.On("GetPlayerByID", mock.Anything, gameID, playerID).
		Return(responsePlayer, nil)

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/players/", nil)
	req.SetPathValue("gameID", gameID.String())
	req.SetPathValue("playerID", playerID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetPlayer(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response player.Player
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Equal(t, response, responsePlayer)

	// Verift mock was called
	mockPlayerService.AssertExpectations(t)
}

func TestHandleGetGamesByID_InvalidGameID(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	playerID := uuid.New()

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/players/", nil)
	req.SetPathValue("gameID", "invalid-id")
	req.SetPathValue("playerID", playerID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetPlayer(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't parse game ID", response["error"])

	// Verify mock was called
	mockPlayerService.AssertNotCalled(t, "GetPlayerByID")
}

func TestHandleGetGamesByID_InvalidPlayerID(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/players/", nil)
	req.SetPathValue("gameID", gameID.String())
	req.SetPathValue("playerID", "invalid-id")
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetPlayer(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Check response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't parse player ID", response["error"])

	// Verify mock was called
	mockPlayerService.AssertNotCalled(t, "GetPlayerByID")
}

func TestHandleGetPlayerByID_ServiceError(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Create test games
	gameID := uuid.New()
	playerID := uuid.New()

	// Set up mock expectations
	mockPlayerService.On("GetPlayerByID", mock.Anything, gameID, playerID).
		Return(player.Player{}, errors.New("error getting player"))

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/players/", nil)
	req.SetPathValue("gameID", gameID.String())
	req.SetPathValue("playerID", playerID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetPlayer(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Check response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "failed to fetch player", response["error"])

	// Verift mock was called
	mockPlayerService.AssertExpectations(t)
}
