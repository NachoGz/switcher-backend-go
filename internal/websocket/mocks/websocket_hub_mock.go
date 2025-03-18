package websocket_mock

import (
	"github.com/NachoGz/switcher-backend-go/internal/websocket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockWebSocketHub struct {
	mock.Mock
}

func (m *MockWebSocketHub) BroadcastToGame(gameID uuid.UUID, messageType string, payload interface{}) {
	m.Called(gameID, messageType, payload)
}

func (m *MockWebSocketHub) BroadcastEvent(gameID uuid.UUID, eventType string) {
	m.Called(gameID, eventType)
}

func (m *MockWebSocketHub) GetClientsInGame(gameID uuid.UUID) int {
	args := m.Called(gameID)
	return args.Int(0)
}

func (m *MockWebSocketHub) RegisterClient(client *websocket.Client) {
	m.Called(client)
}

func (m *MockWebSocketHub) UnregisterClient(client *websocket.Client) {
	m.Called(client)
}

func (m *MockWebSocketHub) BroadcastMessage(message *websocket.BroadcastMessage) {
	m.Called(message)
}
