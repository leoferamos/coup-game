package ws

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client with bidirectional communication capabilities.
type Client struct {
	ID      string             // Unique identifier for the client
	conn    *websocket.Conn    // WebSocket connection
	send    chan []byte        // Buffered channel for outbound messages
	manager *ConnectionManager // Reference to the connection manager
	done    chan struct{}      // Channel to signal client shutdown
}

// NewClient creates a new WebSocket client with the specified parameters.
func NewClient(id string, conn *websocket.Conn, manager *ConnectionManager) *Client {
	if id == "" {
		id = uuid.New().String()
	}

	return &Client{
		ID:      id,
		conn:    conn,
		send:    make(chan []byte, 256),
		manager: manager,
		done:    make(chan struct{}),
	}
}

// StartPumps starts the read and write pumps for the client
func (c *Client) StartPumps() {
	go c.writePump()
	go c.readPump()
}

// writePump handles sending messages to the client
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to current message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-c.done:
			return
		}
	}
}

// readPump handles reading messages from the client
func (c *Client) readPump() {
	defer func() {
		c.manager.RemoveConnection(c.ID)
		c.conn.Close()
		close(c.done)
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse and handle the message
		gameMsg, err := FromJSON(message)
		if err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		c.handleMessage(gameMsg)
	}
}

// handleMessage processes incoming game messages
func (c *Client) handleMessage(msg *GameMessage) {
	log.Printf("Client %s sent message type: %s", c.ID, msg.Type)

	switch msg.Type {
	case Chat:
		// Broadcast chat message to all clients
		c.manager.Broadcast([]byte(fmt.Sprintf("Chat from %s: %v", c.ID, msg.Payload)))
	case PlayerJoin:
		// Handle player join logic
		log.Printf("Player joined: %v", msg.Payload)
	case GameAction:
		// Handle game action
		log.Printf("Game action from %s: %v", c.ID, msg.Payload)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// SendMessage sends a message to the client
func (c *Client) SendMessage(msg *GameMessage) error {
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("client send channel is full")
	}
}

// Close closes the client connection
func (c *Client) Close() {
	close(c.send)
	select {
	case c.done <- struct{}{}:
	default:
	}
}
