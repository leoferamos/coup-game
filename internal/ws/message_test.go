package ws

import (
	"testing"
)

// TDD: Test message protocol structure
func TestGameMessage(t *testing.T) {
	testCases := []struct {
		name     string
		msgType  MessageType
		payload  interface{}
		expected string
	}{
		{
			name:     "Player join message",
			msgType:  PlayerJoin,
			payload:  map[string]string{"playerName": "Alice"},
			expected: `{"type":"player_join","payload":{"playerName":"Alice"}}`,
		},
		{
			name:     "Game state message",
			msgType:  GameState,
			payload:  map[string]int{"turn": 1, "players": 3},
			expected: `{"type":"game_state","payload":{"players":3,"turn":1}}`,
		},
		{
			name:     "Chat message",
			msgType:  Chat,
			payload:  map[string]string{"message": "Hello everyone!"},
			expected: `{"type":"chat","payload":{"message":"Hello everyone!"}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			msg := NewGameMessage(tc.msgType, tc.payload)
			jsonData, err := msg.ToJSON()

			// Assert
			if err != nil {
				t.Errorf("ToJSON() error = %v, want nil", err)
			}
			if string(jsonData) != tc.expected {
				t.Errorf("ToJSON() = %v, want %v", string(jsonData), tc.expected)
			}
		})
	}
}

// TDD: Test message parsing from JSON
func TestFromJSON(t *testing.T) {
	jsonData := `{"type":"chat","payload":{"message":"Hello World"}}`

	msg, err := FromJSON([]byte(jsonData))

	if err != nil {
		t.Errorf("FromJSON() error = %v, want nil", err)
	}

	if msg.Type != Chat {
		t.Errorf("Message type = %v, want %v", msg.Type, Chat)
	}

	// Verify payload structure
	payload, ok := msg.Payload.(map[string]interface{})
	if !ok {
		t.Error("Payload should be map[string]interface{}")
	}

	message, exists := payload["message"]
	if !exists {
		t.Error("Payload should contain 'message' field")
	}

	if message != "Hello World" {
		t.Errorf("Message content = %v, want 'Hello World'", message)
	}
}

// TDD: Test message type validation
func TestMessageTypes(t *testing.T) {
	testCases := []struct {
		msgType  MessageType
		expected string
	}{
		{PlayerJoin, "player_join"},
		{PlayerLeave, "player_leave"},
		{GameState, "game_state"},
		{GameAction, "game_action"},
		{Chat, "chat"},
		{Error, "error"},
	}

	for _, tc := range testCases {
		t.Run(string(tc.msgType), func(t *testing.T) {
			if string(tc.msgType) != tc.expected {
				t.Errorf("MessageType %v = %v, want %v", tc.msgType, string(tc.msgType), tc.expected)
			}
		})
	}
}
