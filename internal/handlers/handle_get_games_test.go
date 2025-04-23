package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NachoGz/switcher-backend-go/internal/game"
	game_mock "github.com/NachoGz/switcher-backend-go/internal/game/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/handlers"
	player_mock "github.com/NachoGz/switcher-backend-go/internal/player/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/websocket"
	websocket_mock "github.com/NachoGz/switcher-backend-go/internal/websocket/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleGetGames_Success(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Create test games
	password := "secret"
	games := []game.Game{
		{
			ID:           uuid.New(),
			Name:         "Test Game 1",
			MaxPlayers:   4,
			MinPlayers:   2,
			PlayersCount: 1,
			IsPrivate:    true,
			Password:     &password,
		},
		{
			ID:           uuid.New(),
			Name:         "Test Game 2",
			MaxPlayers:   4,
			MinPlayers:   3,
			PlayersCount: 2,
			IsPrivate:    false,
		},
	}

	// Set up mock expectations, default parameters
	mockService.On("GetAvailableGames", mock.Anything, 0, 1, 5, "").
		Return(games, 2, nil)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWSHub)

	// Create request without query parameters
	req, _ := http.NewRequest(http.MethodGet, "/games", nil)
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetGames(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify response structure
	assert.Contains(t, response, "games")
	assert.Contains(t, response, "total_pages")
	assert.Equal(t, float64(1), response["total_pages"])

	// Verify games array
	gamesResponse, ok := response["games"].([]interface{}) // Type assersion
	assert.True(t, ok, "games should be an array")
	assert.Equal(t, 2, len(gamesResponse))

	// Verift mock was called
	mockService.AssertExpectations(t)
}

func TestHandleGetGames_withQueryParameters(t *testing.T) {
	// Create mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket.Hub)

	// Create test games
	password := "secret"
	games := []game.Game{
		{
			ID:           uuid.New(),
			Name:         "Test Game 1",
			MaxPlayers:   4,
			MinPlayers:   2,
			PlayersCount: 2,
			IsPrivate:    true,
			Password:     &password,
		},
		{
			ID:           uuid.New(),
			Name:         "Test Game 2",
			MaxPlayers:   4,
			MinPlayers:   3,
			PlayersCount: 2,
			IsPrivate:    false,
		},
	}

	// Set mock expectations with parameters
	mockService.On("GetAvailableGames", mock.Anything, 2, 1, 10, "Game").
		Return(games, 2, nil)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request with query parameters
	req, _ := http.NewRequest(http.MethodGet, "/games?page=1&limit=10&num_players=2&name=Game", nil)
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleGetGames(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify games array
	gamesResponse, ok := response["games"].([]interface{}) // Type assertion
	assert.True(t, ok, "games should be an array")
	assert.Equal(t, 2, len(gamesResponse))

	// Verify mock was called
	mockService.AssertExpectations(t)
}

func TestHandleGetGames_ServiceError(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket.Hub)

	// Setup expectations with error
	mockService.On("GetAvailableGames", mock.Anything, 0, 1, 5, "").
		Return([]game.Game{}, 0, errors.New("database error"))

	// Create handler
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/games", nil)
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetGames(rr, req)

	// Check response
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Error getting games", response["error"])

	// Verify mock was called
	mockService.AssertExpectations(t)
}

func TestHandleGetGames_InvalidPageParameter(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket.Hub)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/games?page=invalid", nil)
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetGames(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid page", response["error"])

	// Verify mock was not called
	mockService.AssertNotCalled(t, "GetAvailableGames")
}

func TestHandleGetGamesByID_Success(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket.Hub)

	// Create test game
	gameID := uuid.New()
	newGame := &game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		PlayersCount: 2,
		IsPrivate:    false,
	}

	// Set up expectations
	mockService.On("GetGameByID", mock.Anything, gameID).
		Return(newGame, nil)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/games/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetGameByID(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response game.Game
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify game response
	assert.Equal(t, gameID, response.ID)
	assert.Equal(t, "Test Game", response.Name)
	assert.Equal(t, 4, response.MaxPlayers)
	assert.Equal(t, 2, response.MinPlayers)
	assert.Equal(t, 2, response.PlayersCount)
	assert.Equal(t, false, response.IsPrivate)

	// Verify mock was called
	mockService.AssertExpectations(t)
}

func TestHandleGetGamesByID_InvalidID(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket_mock.MockWebSocketHub)

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodGet, "/games/", nil)
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
	mockService.AssertNotCalled(t, "GetGameByID")
}

func TestHandleGetGamesByID_ServiceError(t *testing.T) {
	// Setup mock service
	mockService := new(game_mock.MockGameService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockWebsocket := new(websocket_mock.MockWebSocketHub)
	gameID := uuid.New()

	// Setup mock expectations with error
	mockService.On("GetGameByID", mock.Anything, gameID).
		Return((*game.Game)(nil), errors.New("database error"))

	// Create handlers
	handlers := handlers.NewGameHandlers(mockService, mockPlayerService, mockWebsocket)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodGet, "/games/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleGetGameByID(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Error getting game", response["error"])

	// Verify mock was not called
	mockService.AssertExpectations(t)
}
