package ws

import (
	"testing"
)

// TDD: Test NewClient function
func TestNewClient(t *testing.T) {
	manager := NewConnectionManager()

	client := NewClient("test-client-1", nil, manager)
	if client == nil {
		t.Error("NewClient() should return non-nil client")
		return
	}

	if client.ID != "test-client-1" {
		t.Errorf("NewClient() ID = %v, want %v", client.ID, "test-client-1")
	}

	if client.conn != nil {
		t.Error("NewClient() with nil conn should set conn to nil")
	}

	if client.send == nil {
		t.Error("NewClient() should initialize send channel")
	}

	if client.manager != manager {
		t.Error("NewClient() should set manager reference")
	}
}

// TDD: Test Client fields initialization
func TestClientInitialization(t *testing.T) {
	manager := NewConnectionManager()
	client := NewClient("test-client", nil, manager)

	// Test channel is buffered
	select {
	case client.send <- []byte("test"):
		// Should not block due to buffer
	default:
		t.Error("send channel should be buffered")
	}

	// Drain the channel
	<-client.send
}

// TDD: Test client ID generation when empty
func TestClientIDGeneration(t *testing.T) {
	manager := NewConnectionManager()

	// Test creating a client with empty ID - should generate UUID
	client := NewClient("", nil, manager)
	if client == nil {
		t.Error("NewClient() should return non-nil client")
		return
	}

	// Test client ID generation
	if client.ID == "" {
		t.Error("NewClient() should generate non-empty ID when given empty string")
	}

	// Test client manager reference
	if client.manager != manager {
		t.Error("NewClient() should set manager reference")
	}

	// Test that different clients get different IDs
	client2 := NewClient("", nil, manager)
	if client.ID == client2.ID {
		t.Error("NewClient() should generate unique IDs for different clients")
	}
}
