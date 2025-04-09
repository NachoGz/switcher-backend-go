package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	game_mock "github.com/NachoGz/switcher-backend-go/internal/game/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/handlers"
	"github.com/NachoGz/switcher-backend-go/internal/player"
	player_mock "github.com/NachoGz/switcher-backend-go/internal/player/mocks"
	websocket_mock "github.com/NachoGz/switcher-backend-go/internal/websocket/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWinner_Success(t *testing.T) {
	// Setup mock service
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	gameID := uuid.New()
	players := []player.Player{
		{
			ID:          uuid.New(),
			Name:        "Test Player 1",
			Host:        true,
			GameID:      gameID,
			GameStateID: uuid.New(),
			Winner:      true,
		},
		{
			ID:          uuid.New(),
			Name:        "Test Player 2",
			Host:        false,
			GameID:      gameID,
			GameStateID: uuid.New(),
			Winner:      false,
		},
	}

	// Setup mock expectations
	mockPlayerService.On("GetWinner", mock.Anything, gameID).
		Return(&players[0], nil)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockGameService, mockPlayerService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/games/"+gameID.String()+"/winner", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandlerGetWinner(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response player.Player
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Equal(t, response, players[0])

	// Verift mock was called
	mockPlayerService.AssertExpectations(t)
}

func TestHandleGetWinner_InvalidID(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket_mock.MockWebSocketHub)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodGet, "/games/invalid-id/winner", nil)
	req.SetPathValue("gameID", "invalid-ID")
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandlerGetWinner(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't parse game ID", response["error"])

	// Verify mock was not called
	mockService.AssertNotCalled(t, "GetWinner")
}

func TestHandleGetWinner_ServiceError(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket_mock.MockWebSocketHub)
	gameID := uuid.New()

	// Setup mock expectations with error
	mockPlayerService.On("GetWinner", mock.Anything, gameID).
		Return(&player.Player{}, errors.New("database error"))

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodGet, "/games/"+gameID.String()+"/winner", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandlerGetWinner(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't fetch winner", response["error"])

	// Verify mock was not called
	mockService.AssertExpectations(t)
}

func TestHandleGetWinner_NoWinner(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket_mock.MockWebSocketHub)
	gameID := uuid.New()

	// Setup mock expectations with error
	mockPlayerService.On("GetWinner", mock.Anything, gameID).
		Return((*player.Player)(nil), errors.New("database error"))

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodGet, "/games/"+gameID.String()+"/winner", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandlerGetWinner(rr, req)

	// Check response
	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "message")
	assert.Equal(t, "There is no winner", response["message"])

	// Verify mock was not called
	mockService.AssertExpectations(t)
}
