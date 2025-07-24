package lobby

import (
	"fmt"
	"math/rand"
	"time"
)

// Player represents a player in the game
type Player struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Room represents a game room
type Room struct {
	Code       string   `json:"code"`
	Players    []Player `json:"players"`
	MaxPlayers int      `json:"maxPlayers"`
	CreatedAt  time.Time `json:"createdAt"`
}

// CreateRoom creates a new room with a 4-digit numeric code
func CreateRoom() *Room {
	return &Room{
		Code:       generateRoomCode(),
		Players:    make([]Player, 0),
		MaxPlayers: 10, // Maximum players for Coup expansion
		CreatedAt:  time.Now(),
	}
}

// generateRoomCode generates a 4-digit numeric code
func generateRoomCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(10000) // 0-9999
	return fmt.Sprintf("%04d", code) // Ensure 4 digits with leading zeros
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(player Player) error {
	// Check if room is full
	if len(r.Players) >= r.MaxPlayers {
		return fmt.Errorf("room is full, maximum %d players allowed", r.MaxPlayers)
	}
	
	// Check for duplicate player ID
	for _, existingPlayer := range r.Players {
		if existingPlayer.ID == player.ID {
			return fmt.Errorf("player with ID %s already exists in room", player.ID)
		}
	}
	
	// Add player to room
	r.Players = append(r.Players, player)
	return nil
}

// RemovePlayer removes a player from the room by ID
func (r *Room) RemovePlayer(playerID string) error {
	for i, player := range r.Players {
		if player.ID == playerID {
			// Remove player from slice
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("player with ID %s not found in room", playerID)
}

// IsReadyToStart checks if the room has enough players to start a game
func (r *Room) IsReadyToStart() bool {
	// Coup requires at least 3 players
	return len(r.Players) >= 3
}

// CalculateDeckSize calculates the total deck size based on player count
func CalculateDeckSize(playerCount int) int {
	cardsPerInfluence := CalculateCardsPerInfluence(playerCount)
	// 5 influences in Coup: Duke, Assassin, Captain, Ambassador, Contessa
	return cardsPerInfluence * 5
}

// CalculateCardsPerInfluence calculates cards per influence based on player count
func CalculateCardsPerInfluence(playerCount int) int {
	if playerCount >= 3 && playerCount <= 6 {
		return 3 // 3 of each for 3-6 players
	} else if playerCount >= 7 && playerCount <= 8 {
		return 4 // 4 of each for 7-8 players
	} else if playerCount >= 9 && playerCount <= 10 {
		return 5 // 5 of each for 9-10 players
	}
	
	// Default to 3 for invalid player counts
	return 3
}
