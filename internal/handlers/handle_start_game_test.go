package handlers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	board_mock "github.com/NachoGz/switcher-backend-go/internal/board/mocks"
	figureCard_mock "github.com/NachoGz/switcher-backend-go/internal/figureCard/mocks"
	gameState "github.com/NachoGz/switcher-backend-go/internal/game_state"
	gameState_mock "github.com/NachoGz/switcher-backend-go/internal/game_state/mocks"
	"github.com/NachoGz/switcher-backend-go/internal/handlers"
	movementCard_mock "github.com/NachoGz/switcher-backend-go/internal/movementCard/mocks"
	websocket_mock "github.com/NachoGz/switcher-backend-go/internal/websocket/mocks"

	"github.com/NachoGz/switcher-backend-go/internal/player"
	player_mock "github.com/NachoGz/switcher-backend-go/internal/player/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleStartGame_Success(t *testing.T) {
	// Setup mock
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()

	responsePlayers := []player.Player{
		{
			ID:          uuid.New(),
			Name:        "Test Player 1",
			Host:        true,
			GameID:      gameID,
			GameStateID: gameStateID},
		{
			ID:          uuid.New(),
			Name:        "Test Player 2",
			Host:        true,
			GameID:      gameID,
			GameStateID: gameStateID},
	}

	// Setup expectations
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(nil)

	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return(responsePlayers, nil)

	mockPlayerService.On("AssignRandomTurns", mock.Anything, responsePlayers).
		Return(responsePlayers[0].ID, nil)

	mockGameStateService.On("UpdateCurrentPlayer", mock.Anything, gameID, responsePlayers[0].ID).
		Return(nil)

	mockBoardService.On("ConfigureBoard", mock.Anything, gameID).
		Return(nil)

	mockMovementCardService.On("CreateMovementCardDeck", mock.Anything, gameID).
		Return(nil)

	mockFigureCardService.On("CreateFigureCardDeck", mock.Anything, gameID).
		Return(nil)

	mockWSHub.On("BroadcastEvent", uuid.Nil, "GAMES_LIST_UPDATE").
		Return()

	mockWSHub.On("BroadcastEvent", uuid.Nil, fmt.Sprintf("%s:GAME_STARTED", gameID)).
		Return()

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService, mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verify response body
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Equal(t, "Game started successfully", response["message"])

	// Verify mocks are called
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockBoardService.AssertExpectations(t)
	mockMovementCardService.AssertExpectations(t)
	mockFigureCardService.AssertExpectations(t)
	mockWSHub.AssertExpectations(t)
}

func TestHandleStartGame_InvalidGameID(t *testing.T) {
	// Setup mock
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Create handlers with mock service
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService, mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create invalid request body
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", "invalid-id")
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't parse game ID", response["error"])

	// Ensure services are not called
	mockGameStateService.AssertNotCalled(t, "UpdateGameState")
	mockPlayerService.AssertNotCalled(t, "GetPlayersInGame")
	mockPlayerService.AssertNotCalled(t, "AssignRandomTurns")
	mockBoardService.AssertNotCalled(t, "ConfigureBoard")
	mockMovementCardService.AssertNotCalled(t, "CreateMovementCardDeck")
	mockFigureCardService.AssertNotCalled(t, "CreateFigureCardDeck")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleStartGame_UpdateGameStateError(t *testing.T) {
	// Setup mocks
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()

	// Mock error
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(errors.New("database error"))

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService,
		mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify error message
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "Error updating game state", response["error"])

	// Verify only the necessary services were called
	mockGameStateService.AssertExpectations(t)
	mockGameStateService.AssertNotCalled(t, "UpdateGameState")
	mockPlayerService.AssertNotCalled(t, "GetPlayersInGame")
	mockPlayerService.AssertNotCalled(t, "AssignRandomTurns")
	mockBoardService.AssertNotCalled(t, "ConfigureBoard")
	mockMovementCardService.AssertNotCalled(t, "CreateMovementCardDeck")
	mockFigureCardService.AssertNotCalled(t, "CreateFigureCardDeck")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleStartGame_GetPlayersError(t *testing.T) {
	// Setup mocks
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()

	// Mock responses
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(nil)
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return([]player.Player{}, errors.New("database error when fetching players"))

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService,
		mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify error message
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "Error fetching players", response["error"])

	// Verify only the necessary services were called
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertNotCalled(t, "AssignRandomTurns")
	mockGameStateService.AssertNotCalled(t, "UpdateCurrentPlayer")
	mockBoardService.AssertNotCalled(t, "ConfigureBoard")
	mockMovementCardService.AssertNotCalled(t, "CreateMovementCardDeck")
	mockFigureCardService.AssertNotCalled(t, "CreateFigureCardDeck")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleStartGame_AssignRandomTurnError(t *testing.T) {
	// Setup mocks
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()

	// Mock responses
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(nil)
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return([]player.Player{}, nil)

	mockPlayerService.On("AssignRandomTurns", mock.Anything, []player.Player{}).
		Return(uuid.Nil, errors.New("there are no players"))

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService,
		mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify error message
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "Error setting turns")

	// Verify only the necessary services were called
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockGameStateService.AssertNotCalled(t, "UpdateCurrentPlayer")
	mockBoardService.AssertNotCalled(t, "ConfigureBoard")
	mockMovementCardService.AssertNotCalled(t, "CreateMovementCardDeck")
	mockFigureCardService.AssertNotCalled(t, "CreateFigureCardDeck")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleStartGame_UpdateCurrentPlayerError(t *testing.T) {
	// Setup mocks
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
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

	// Mock responses
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(nil)
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return(responsePlayers, nil)
	mockPlayerService.On("AssignRandomTurns", mock.Anything, responsePlayers).
		Return(responsePlayers[0].ID, nil)
	mockGameStateService.On("UpdateCurrentPlayer", mock.Anything, gameID, responsePlayers[0].ID).
		Return(errors.New("error updating current player"))

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService,
		mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify error message
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "Error updating current player", response["error"])

	// Verify only the necessary services were called
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockBoardService.AssertNotCalled(t, "ConfigureBoard")
	mockMovementCardService.AssertNotCalled(t, "CreateMovementCardDeck")
	mockFigureCardService.AssertNotCalled(t, "CreateFigureCardDeck")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleStartGame_ConfigureBoardError(t *testing.T) {
	// Setup mocks
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
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
	}

	// Mock responses
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(nil)
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return(responsePlayers, nil)
	mockPlayerService.On("AssignRandomTurns", mock.Anything, responsePlayers).
		Return(responsePlayers[0].ID, nil)
	mockGameStateService.On("UpdateCurrentPlayer", mock.Anything, gameID, responsePlayers[0].ID).
		Return(nil)
	mockBoardService.On("ConfigureBoard", mock.Anything, gameID).
		Return(errors.New("error configuring board"))

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService,
		mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify error message
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "Error configuring board", response["error"])

	// Verify only the necessary services were called
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockBoardService.AssertExpectations(t)
	mockMovementCardService.AssertNotCalled(t, "CreateMovementCardDeck")
	mockFigureCardService.AssertNotCalled(t, "CreateFigureCardDeck")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleStartGame_CreateMovementCardDeckError(t *testing.T) {
	// Setup mocks
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
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
	}

	// Mock responses
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(nil)
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return(responsePlayers, nil)
	mockPlayerService.On("AssignRandomTurns", mock.Anything, responsePlayers).
		Return(responsePlayers[0].ID, nil)
	mockGameStateService.On("UpdateCurrentPlayer", mock.Anything, gameID, responsePlayers[0].ID).
		Return(nil)
	mockBoardService.On("ConfigureBoard", mock.Anything, gameID).
		Return(nil)

	mockMovementCardService.On("CreateMovementCardDeck", mock.Anything, gameID).
		Return(errors.New("error creating movement card deck"))

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService,
		mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify error message
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "Error creating movement card deck", response["error"])

	// Verify only the necessary services were called
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockBoardService.AssertExpectations(t)
	mockMovementCardService.AssertExpectations(t)
	mockFigureCardService.AssertNotCalled(t, "CreateFigureCardDeck")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleStartGame_CreateFigureCardDeckError(t *testing.T) {
	// Setup mocks
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockPlayerService := new(player_mock.MockPlayerService)
	mockBoardService := new(board_mock.MockBoardService)
	mockMovementCardService := new(movementCard_mock.MockMovementCardService)
	mockFigureCardService := new(figureCard_mock.MockFigureCardService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
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
	}

	// Mock responses
	mockGameStateService.On("UpdateGameState", mock.Anything, gameID, gameState.PLAYING).
		Return(nil)
	mockPlayerService.On("GetPlayersInGame", mock.Anything, gameID).
		Return(responsePlayers, nil)
	mockPlayerService.On("AssignRandomTurns", mock.Anything, responsePlayers).
		Return(responsePlayers[0].ID, nil)
	mockGameStateService.On("UpdateCurrentPlayer", mock.Anything, gameID, responsePlayers[0].ID).
		Return(nil)
	mockBoardService.On("ConfigureBoard", mock.Anything, gameID).
		Return(nil)

	mockMovementCardService.On("CreateMovementCardDeck", mock.Anything, gameID).
		Return(nil)

	mockFigureCardService.On("CreateFigureCardDeck", mock.Anything, gameID).
		Return(errors.New("error creating figure card deck"))

	// Create handlers
	handlers := handlers.NewGameStateHandlers(mockGameStateService, mockPlayerService,
		mockBoardService, mockMovementCardService, mockFigureCardService, mockWSHub)

	// Create request
	req, _ := http.NewRequest(http.MethodPatch, "/games/start/", nil)
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleStartGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verify error message
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "error")
	assert.Equal(t, "Error creating figure card deck", response["error"])

	// Verify the services are called
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockBoardService.AssertExpectations(t)
	mockMovementCardService.AssertExpectations(t)
	mockFigureCardService.AssertExpectations(t)
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}
