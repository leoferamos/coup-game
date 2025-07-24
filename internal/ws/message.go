package ws

import (
	"encoding/json"
	"time"
)

// MessageType represents different types of game messages
type MessageType string

const (
	PlayerJoin  MessageType = "player_join"
	PlayerLeave MessageType = "player_leave"
	GameState   MessageType = "game_state"
	GameAction  MessageType = "game_action"
	Chat        MessageType = "chat"
	Error       MessageType = "error"
)

// GameMessage represents a structured message in the game protocol
type GameMessage struct {
	Type      MessageType `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
}

// NewGameMessage creates a new game message with timestamp
func NewGameMessage(msgType MessageType, payload interface{}) *GameMessage {
	return &GameMessage{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now(),
	}
}

// ToJSON converts the message to JSON bytes
func (gm *GameMessage) ToJSON() ([]byte, error) {
	msg := struct {
		Type    MessageType `json:"type"`
		Payload interface{} `json:"payload"`
	}{
		Type:    gm.Type,
		Payload: gm.Payload,
	}
	return json.Marshal(msg)
}

// FromJSON creates a GameMessage from JSON bytes
func FromJSON(data []byte) (*GameMessage, error) {
	var msg GameMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
