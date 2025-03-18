package websocket

import "github.com/google/uuid"

// WebSocketHub defines methods for broadcasting and managing connections.
type WebSocketHub interface {
	BroadcastToGame(gameID uuid.UUID, messageType string, payload interface{})
	BroadcastEvent(gameID uuid.UUID, eventType string)
	GetClientsInGame(gameID uuid.UUID) int
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	BroadcastMessage(message *BroadcastMessage)
}

// Ensure Hub implements WebSocketHub
var _ WebSocketHub = (*Hub)(nil)
