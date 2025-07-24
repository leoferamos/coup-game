package ws

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager manages all active WebSocket connections in a thread-safe manner.
type ConnectionManager struct {
	connections map[string]*Client // Active client connections indexed by ID
	mu          sync.RWMutex       // Read-write mutex for thread-safe operations
}

// NewConnectionManager creates a new ConnectionManager instance.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Client),
	}
}

// GetConnectionCount returns the number of active connections
func (cm *ConnectionManager) GetConnectionCount() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.connections)
}

// AddConnection adds a new connection to the manager (legacy method)
func (cm *ConnectionManager) AddConnection(id string, conn *websocket.Conn) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if id == "" {
		return fmt.Errorf("connection ID cannot be empty")
	}

	if _, exists := cm.connections[id]; exists {
		return fmt.Errorf("connection with ID %s already exists", id)
	}

	client := NewClient(id, conn, cm)
	cm.connections[id] = client

	return nil
}

// AddClient adds a new client to the manager
func (cm *ConnectionManager) AddClient(client *Client) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client.ID == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	if _, exists := cm.connections[client.ID]; exists {
		return fmt.Errorf("client with ID %s already exists", client.ID)
	}

	cm.connections[client.ID] = client
	return nil
}

// RemoveConnection removes a connection from the manager
func (cm *ConnectionManager) RemoveConnection(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exists := cm.connections[id]
	if !exists {
		return fmt.Errorf("connection with ID %s not found", id)
	}

	// Close the client properly
	client.Close()
	delete(cm.connections, id)

	return nil
}

// GetConnections returns a copy of connection IDs for safe iteration
func (cm *ConnectionManager) GetConnections() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	connections := make([]string, 0, len(cm.connections))
	for id := range cm.connections {
		connections = append(connections, id)
	}

	return connections
}

// Broadcast sends a message to all active connections
func (cm *ConnectionManager) Broadcast(message []byte) error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	for _, client := range cm.connections {
		select {
		case client.send <- message:
		default:
			// Connection is blocked, skip it
			log.Printf("Skipping blocked connection: %s", client.ID)
		}
	}

	return nil
}

// BroadcastMessage broadcasts a structured game message to all clients
func (cm *ConnectionManager) BroadcastMessage(msg *GameMessage) error {
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}
	return cm.Broadcast(data)
}
