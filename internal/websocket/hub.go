package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Client struct {
	Server   WebSocketHub
	Conn     *Connection
	Send     chan []byte
	GameID   uuid.UUID
	PlayerID uuid.UUID
}

// Message represents a structured message for WebSocket communication
type Message struct {
	Type    string      `json:"type"`
	GameID  string      `json:"game_id,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by game ID
	clients map[uuid.UUID]map[*Client]bool

	// Register requests from clients
	Register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast message to specific game
	broadcast chan *BroadcastMessage

	// Mutex to protect concurrent access to the clients map
	mu sync.Mutex
}

// BroadcastMessage contains the message data and target game
type BroadcastMessage struct {
	GameID  uuid.UUID
	Message []byte
}

// NewServer creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*Client]bool),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage),
	}
}

// Run starts the Hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.clients[client.GameID]; !ok {
				h.clients[client.GameID] = make(map[*Client]bool)
			}
			h.clients[client.GameID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered for game %s, total clients: %d",
				client.GameID, len(h.clients[client.GameID]))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.GameID]; ok {
				if _, ok := h.clients[client.GameID][client]; ok {
					delete(h.clients[client.GameID], client)
					close(client.Send)
					log.Printf("Client unregistered from game %s", client.GameID)

					// If no clients left in the game, clean up
					if len(h.clients[client.GameID]) == 0 {
						delete(h.clients, client.GameID)
						log.Printf("No clients left in game %s, removing game", client.GameID)
					}
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			if clients, ok := h.clients[message.GameID]; ok {
				for client := range clients {
					select {
					case client.Send <- message.Message:
					default:
						close(client.Send)
						delete(h.clients[message.GameID], client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastToGame sends a JSON message to all clients in a specific game
func (h *Hub) BroadcastToGame(gameID uuid.UUID, messageType string, payload interface{}) {
	// Create the message structure
	message := Message{
		Type:    messageType,
		GameID:  gameID.String(),
		Payload: payload,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message to JSON: %v", err)
		return
	}

	log.Printf("Broadcasting to game %s: %s", gameID, string(jsonData))

	// Send through the broadcast channel
	h.broadcast <- &BroadcastMessage{
		GameID:  gameID,
		Message: jsonData,
	}
}

// BroadcastEvent sends a simple event message (no payload) to all clients in a game
func (h *Hub) BroadcastEvent(gameID uuid.UUID, eventType string) {
	h.BroadcastToGame(gameID, eventType, nil)
}

// GetClientsInGame returns the number of clients connected to a specific game
func (h *Hub) GetClientsInGame(gameID uuid.UUID) int {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.clients[gameID]; ok {
		return len(clients)
	}
	return 0
}

// RegisterClient registers a client with the hub
func (h *Hub) RegisterClient(client *Client) {
	h.Register <- client
}

// UnregisterClient unregisters a client from the hub
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastMessage sends a message through the broadcast channel
func (h *Hub) BroadcastMessage(message *BroadcastMessage) {
	h.broadcast <- message
}
