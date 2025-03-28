package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	game_mock "github.com/NachoGz/switcher-backend-go/internal/game/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/handlers"
	"github.com/NachoGz/switcher-backend-go/internal/websocket"
	websocket_mock "github.com/NachoGz/switcher-backend-go/internal/websocket/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleDeleteGame_Success(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Create test game
	gameID := uuid.New()

	// Setup expectations
	mockService.On("DeleteGame", mock.Anything, gameID).
		Return(nil)

	mockWSHub.On("BroadcastEvent", uuid.Nil, "GAMES_LIST_UPDATE").
		Return()

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodDelete, "/games/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleDeleteGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusNoContent, rr.Code)

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify  response
	assert.Contains(t, response, "message")
	assert.Equal(t, "Game deleted successfully", response["message"])

	// Verify mock was called
	mockService.AssertExpectations(t)
	mockWSHub.AssertExpectations(t)
}

func TestHandleDeleteGame_InvalidID(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockWebsocket := new(websocket.Hub)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockWebsocket)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodDelete, "/games/", nil)
	req.SetPathValue("gameID", "invalid-ID")
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetGameByID(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't parse game ID", response["error"])

	// Verify mock was not called
	mockService.AssertNotCalled(t, "DeleteGame")
}

func TestHandleDeleteGame_ServiceError(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockWebsocket := new(websocket.Hub)
	gameID := uuid.New()

	// Setup mock expectations with error
	mockService.On("DeleteGame", mock.Anything, gameID).
		Return(errors.New("database error"))

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockWebsocket)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodDelete, "/games/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleDeleteGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't delete game", response["error"])

	// Verify mock was not called
	mockService.AssertExpectations(t)
}
