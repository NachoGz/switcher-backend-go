package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins for development
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Connection represents a websocket connection
type Connection struct {
	ws *websocket.Conn
}

// NewConnection creates a new connection from an http request
func NewConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &Connection{ws: ws}, nil
}

// ReadMessage reads a message from the websocket connection
func (c *Connection) ReadMessage() (int, []byte, error) {
	return c.ws.ReadMessage()
}

// WriteMessage writes a message to the websocket connection
func (c *Connection) WriteMessage(messageType int, data []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(messageType, data)
}

// Close closes the websocket connection
func (c *Connection) Close() error {
	return c.ws.Close()
}

// ReadPump pumps messages from the websocket to the server
func (c *Client) ReadPump() {
	defer func() {
		c.Server.UnregisterClient(c)
		c.Conn.Close()
	}()

	c.Conn.ws.SetReadLimit(maxMessageSize)
	c.Conn.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.ws.SetPongHandler(func(string) error {
		c.Conn.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		if eventType, ok := msg["type"].(string); ok {
			switch eventType {
			case "GAMES_LIST_UPDATE":
				log.Println("There was an update in the games list")
			default:
				log.Printf("Unknown event type: %s", eventType)
			}
		}

		c.Server.BroadcastMessage(&BroadcastMessage{
			GameID:  c.GameID,
			Message: message,
		})
	}
}

// WritePump pumps messages from the server to the websocket connection
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// The server closed the channel
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
