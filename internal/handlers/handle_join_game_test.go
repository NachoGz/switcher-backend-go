package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/NachoGz/switcher-backend-go/internal/utils"
	websocket_mock "github.com/NachoGz/switcher-backend-go/internal/websocket/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleJoinGamePublic_Success(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

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

	mockWSHub.On("BroadcastEvent", uuid.Nil, "GAMES_LIST_UPDATE").
		Return()

	mockWSHub.On("BroadcastEvent", uuid.Nil, fmt.Sprintf("%s:GAME_INFO_UPDATE", gameID)).
		Return()

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

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
	mockWSHub.AssertExpectations(t)
}

func TestHandleJoinGamePrivate_Success(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	playerID := uuid.New()
	password := "password"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Create expected response objects
	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
		Password:     &hashedPassword,
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

	mockWSHub.On("BroadcastEvent", uuid.Nil, "GAMES_LIST_UPDATE").
		Return()

	mockWSHub.On("BroadcastEvent", uuid.Nil, fmt.Sprintf("%s:GAME_INFO_UPDATE", gameID)).
		Return()

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    password,
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
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Contains(t, response, "message")
	assert.Equal(t, "Joined game successfully", response["message"])

	// Verify mocks are called
	mockPlayerService.AssertExpectations(t)
	mockGameService.AssertExpectations(t)
	mockGameStateService.AssertExpectations(t)
	mockWSHub.AssertExpectations(t)
}

func TestHandleJoinGamePrivate_IncorrectPassword(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	password := "password"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Create expected response objects
	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
		Password:     &hashedPassword,
	}

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&responseGame, nil)

	mockPlayerService.On("CountPlayers", mock.Anything, gameID).
		Return(responseGame.PlayersCount, nil)

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "wrongPassword",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Verify response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Incorrect password", response["error"])

	// Verify only the necessary services were called
	mockPlayerService.AssertExpectations(t)
	mockGameService.AssertExpectations(t)
	mockGameStateService.AssertExpectations(t)
	mockGameStateService.AssertNotCalled(t, "GetGameStateByGameID")
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")

}

func TestHandleJoinGamePrivate_NoPasswordGiven(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	password := "password"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	// Create expected response objects
	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
		Password:     &hashedPassword,
	}

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&responseGame, nil)

	mockPlayerService.On("CountPlayers", mock.Anything, gameID).
		Return(responseGame.PlayersCount, nil)

	// Create handlers
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Verify response body
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Password required for private games", response["error"])

	// Verify only the necessary services were called
	mockPlayerService.AssertExpectations(t)
	mockGameService.AssertExpectations(t)
	mockGameStateService.AssertExpectations(t)
	mockGameStateService.AssertNotCalled(t, "GetGameStateByGameID")
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleJoinGame_InvalidGameID(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	// Create handlers with mock service
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "password",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	// Create request with invalid gameID
	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", "invalid-id")
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't parse game ID", response["error"])

	// Ensure services are not called
	mockGameService.AssertNotCalled(t, "GetGameByID")
	mockPlayerService.AssertNotCalled(t, "CountPlayers")
	mockGameStateService.AssertNotCalled(t, "GetGameStateByGameID")
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleJoinGame_InvalidRequestBody(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()

	// Create handlers with mock service
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	reqBodyBytes := []byte(`{invalid json}`)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid request body", response["error"])

	// Ensure services are not called
	mockGameService.AssertNotCalled(t, "GetGameByID")
	mockPlayerService.AssertNotCalled(t, "CountPlayers")
	mockGameStateService.AssertNotCalled(t, "GetGameStateByGameID")
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleJoinGame_GameNotFound(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&game.Game{}, errors.New("game not found"))

	// Create handlers with mock service
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "password",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, fmt.Sprintf("Couldn't get game with ID: %s", gameID), response["error"])

	// Ensure services are not called
	mockGameService.AssertExpectations(t)
	mockPlayerService.AssertNotCalled(t, "CountPlayers")
	mockGameStateService.AssertNotCalled(t, "GetGameStateByGameID")
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleJoinGame_NoPlayersInGame(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	password := "password"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
		Password:     &hashedPassword,
	}

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&responseGame, nil)

	mockPlayerService.On("CountPlayers", mock.Anything, gameID).
		Return(0, errors.New("no players in game"))

	// Create handlers with mock service
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "password",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, fmt.Sprintf("Couldn't get amount of players in game with ID: %s", gameID), response["error"])

	// Ensure services are not called
	mockGameService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockGameStateService.AssertNotCalled(t, "GetGameStateByGameID")
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleJoinGame_FullGame(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	password := "password"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
		Password:     &hashedPassword,
	}

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&responseGame, nil)

	mockPlayerService.On("CountPlayers", mock.Anything, gameID).
		Return(4, nil)

	// Create handlers with mock service
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "password",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusForbidden, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "The game is full", response["error"])

	// Ensure services are not called
	mockGameService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockGameStateService.AssertNotCalled(t, "GetGameStateByGameID")
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleJoinGame_GameStateError(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	password := "password"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
		Password:     &hashedPassword,
	}

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&responseGame, nil)

	mockPlayerService.On("CountPlayers", mock.Anything, gameID).
		Return(responseGame.PlayersCount, nil)

	mockGameStateService.On("GetGameStateByGameID", mock.Anything, gameID).
		Return(&gameState.GameState{}, errors.New("game state not found"))

	// Create handlers with mock service
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "password",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, fmt.Sprintf("Couldn't get game state for game %s", gameID), response["error"])

	// Ensure services are not called
	mockGameService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertNotCalled(t, "CreatePlayer")
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}

func TestHandleJoinGame_CreatePlayerError(t *testing.T) {
	// Setup mock
	mockPlayerService := new(player_mock.MockPlayerService)
	mockGameService := new(game_mock.MockGameService)
	mockGameStateService := new(gameState_mock.MockGameStateService)
	mockWSHub := new(websocket_mock.MockWebSocketHub)

	// Test data
	gameID := uuid.New()
	gameStateID := uuid.New()
	password := "password"

	hashedPassword, err := utils.HashPassword(password)
	assert.NoError(t, err)

	responseGame := game.Game{
		ID:           gameID,
		Name:         "Test Game",
		MaxPlayers:   4,
		MinPlayers:   2,
		IsPrivate:    true,
		PlayersCount: 1,
		Password:     &hashedPassword,
	}

	requestPlayer := player.Player{
		GameID:      gameID,
		GameStateID: gameStateID,
		Name:        "Test Player",
		Host:        false,
	}

	responseGameState := gameState.GameState{
		ID:     gameStateID,
		State:  gameState.WAITING,
		GameID: gameID,
	}

	mockGameService.On("GetGameByID", mock.Anything, gameID).
		Return(&responseGame, nil)

	mockPlayerService.On("CountPlayers", mock.Anything, gameID).
		Return(responseGame.PlayersCount, nil)

	mockGameStateService.On("GetGameStateByGameID", mock.Anything, gameID).
		Return(&responseGameState, nil)

	mockPlayerService.On("CreatePlayer", mock.Anything, requestPlayer).
		Return(&player.Player{}, errors.New("Couldn't create player"))

	// Create handlers with mock service
	handlers := handlers.NewPlayerHandlers(mockPlayerService, mockGameService, mockGameStateService, mockWSHub)

	// Create request body
	requestBody := map[string]interface{}{
		"player_name": requestPlayer.Name,
		"password":    "password",
	}
	reqBodyBytes, _ := json.Marshal(requestBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/games/join/", bytes.NewReader(reqBodyBytes))
	req.SetPathValue("gameID", gameID.String())
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.HandleJoinGame(rr, req)

	// Check response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Verify error message
	assert.Contains(t, response, "error")
	assert.Equal(t, "Couldn't create player", response["error"])

	// Ensure services are not called
	mockGameService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockGameStateService.AssertExpectations(t)
	mockPlayerService.AssertExpectations(t)
	mockWSHub.AssertNotCalled(t, "BroadcastEvent")
}
