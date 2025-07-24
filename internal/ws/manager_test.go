package ws

import (
	"fmt"
	"testing"
	"time"
)

// TDD: Test ConnectionManager basic functionality
func TestConnectionManager(t *testing.T) {
	// Test that we can track active connections
	manager := NewConnectionManager()

	if manager == nil {
		t.Error("NewConnectionManager() should return non-nil manager")
	}

	// Test initial state
	count := manager.GetConnectionCount()
	if count != 0 {
		t.Errorf("GetConnectionCount() = %v, want 0 for new manager", count)
	}
}

// TDD: Test message broadcasting
func TestBroadcastMessage(t *testing.T) {
	manager := NewConnectionManager()

	// Test broadcasting to empty connection pool
	err := manager.Broadcast([]byte("test message"))
	if err != nil {
		t.Errorf("Broadcast() to empty pool should not return error, got: %v", err)
	}
}

// TDD: Test connection management interface
func TestConnectionManagerInterface(t *testing.T) {
	manager := NewConnectionManager()

	// Test adding a connection
	connID := "test-conn-1"
	err := manager.AddConnection(connID, nil)
	if err != nil {
		t.Errorf("AddConnection() error = %v, want nil", err)
	}

	// Test connection count
	count := manager.GetConnectionCount()
	if count != 1 {
		t.Errorf("GetConnectionCount() = %v, want 1", count)
	}

	// Test removing connection
	err = manager.RemoveConnection(connID)
	if err != nil {
		t.Errorf("RemoveConnection() error = %v, want nil", err)
	}

	// Test connection count after removal
	count = manager.GetConnectionCount()
	if count != 0 {
		t.Errorf("GetConnectionCount() after removal = %v, want 0", count)
	}
}

// TDD: Test for production-safe error handling
func TestGetConnectionsSafe(t *testing.T) {
	manager := NewConnectionManager()

	// Test getting connections safely
	connections := manager.GetConnections()
	if connections == nil {
		t.Error("GetConnections() should return non-nil slice")
	}
	if len(connections) != 0 {
		t.Errorf("GetConnections() length = %v, want 0 for empty manager", len(connections))
	}
}

// TDD: Test concurrent connection management
func TestConcurrentConnections(t *testing.T) {
	manager := NewConnectionManager()
	const numConnections = 10

	// Add connections concurrently
	for i := 0; i < numConnections; i++ {
		go func(id int) {
			connID := fmt.Sprintf("conn-%d", id)
			err := manager.AddConnection(connID, nil)
			if err != nil {
				t.Errorf("AddConnection(%s) error = %v, want nil", connID, err)
			}
		}(i)
	}

	// Give time for goroutines to complete
	time.Sleep(100 * time.Millisecond)

	// Check final count
	count := manager.GetConnectionCount()
	if count != numConnections {
		t.Errorf("GetConnectionCount() = %v, want %v", count, numConnections)
	}
}

// TDD: Test AddClient method
func TestAddClient(t *testing.T) {
	manager := NewConnectionManager()
	client := NewClient("test-client", nil, manager)

	err := manager.AddClient(client)
	if err != nil {
		t.Errorf("AddClient() error = %v, want nil", err)
	}

	count := manager.GetConnectionCount()
	if count != 1 {
		t.Errorf("GetConnectionCount() after AddClient = %v, want 1", count)
	}

	// Test adding duplicate client
	err = manager.AddClient(client)
	if err == nil {
		t.Error("AddClient() with duplicate ID should return error")
	}
}

// TDD: Test BroadcastMessage with structured messages
func TestBroadcastMessageStructured(t *testing.T) {
	manager := NewConnectionManager()

	// Create a test message
	msg := NewGameMessage(PlayerJoin, map[string]string{
		"playerId": "test-player",
		"message":  "Test message",
	})

	err := manager.BroadcastMessage(msg)
	if err != nil {
		t.Errorf("BroadcastMessage() error = %v, want nil", err)
	}
}
