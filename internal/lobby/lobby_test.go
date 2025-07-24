package lobby

import (
	"fmt"
	"testing"
)

// TDD: Test room creation with 4-digit numeric code
func TestCreateRoom(t *testing.T) {
	room := CreateRoom()

	if room == nil {
		t.Error("CreateRoom() should return non-nil room")
		return
	}

	// Test room code is 4 digits
	if len(room.Code) != 4 {
		t.Errorf("Room code length = %v, want 4", len(room.Code))
	}

	// Test room code contains only digits
	for _, char := range room.Code {
		if char < '0' || char > '9' {
			t.Errorf("Room code contains non-digit character: %c", char)
		}
	}

	// Test room starts empty
	if len(room.Players) != 0 {
		t.Errorf("New room should start with 0 players, got %v", len(room.Players))
	}

	// Test room has correct max players (10 for full expansion)
	if room.MaxPlayers != 10 {
		t.Errorf("Room max players = %v, want 10", room.MaxPlayers)
	}
}

// TDD: Test player joining room with valid code
func TestJoinRoom(t *testing.T) {
	room := CreateRoom()

	player := Player{
		ID:   "player-1",
		Name: "Alice",
	}

	// Test successful join
	err := room.AddPlayer(player)
	if err != nil {
		t.Errorf("AddPlayer() error = %v, want nil", err)
	}

	// Test player count
	if len(room.Players) != 1 {
		t.Errorf("Room players count = %v, want 1", len(room.Players))
	}

	// Test player is in room
	if room.Players[0].ID != player.ID {
		t.Errorf("Player ID = %v, want %v", room.Players[0].ID, player.ID)
	}
}

// TDD: Test room capacity limits
func TestRoomCapacity(t *testing.T) {
	room := CreateRoom()

	// Fill room to capacity (10 players)
	for i := 0; i < 10; i++ {
		player := Player{
			ID:   fmt.Sprintf("player-%d", i),
			Name: fmt.Sprintf("Player%d", i),
		}

		err := room.AddPlayer(player)
		if err != nil {
			t.Errorf("AddPlayer(%d) error = %v, want nil", i, err)
		}
	}

	// Test room is full
	if len(room.Players) != 10 {
		t.Errorf("Room players count = %v, want 10", len(room.Players))
	}

	// Test adding 11th player fails
	extraPlayer := Player{ID: "player-11", Name: "Extra"}
	err := room.AddPlayer(extraPlayer)
	if err == nil {
		t.Error("AddPlayer() to full room should return error")
	}
}

// TDD: Test duplicate player prevention
func TestDuplicatePlayerPrevention(t *testing.T) {
	room := CreateRoom()

	player := Player{
		ID:   "player-1",
		Name: "Alice",
	}

	// Add player first time
	err := room.AddPlayer(player)
	if err != nil {
		t.Errorf("First AddPlayer() error = %v, want nil", err)
	}

	// Try to add same player again
	err = room.AddPlayer(player)
	if err == nil {
		t.Error("AddPlayer() with duplicate ID should return error")
	}

	// Verify player count is still 1
	if len(room.Players) != 1 {
		t.Errorf("Room players count = %v, want 1", len(room.Players))
	}
}

// TDD: Test player removal
func TestRemovePlayer(t *testing.T) {
	room := CreateRoom()

	player := Player{
		ID:   "player-1",
		Name: "Alice",
	}

	// Add player
	room.AddPlayer(player)

	// Remove player
	err := room.RemovePlayer("player-1")
	if err != nil {
		t.Errorf("RemovePlayer() error = %v, want nil", err)
	}

	// Test player count
	if len(room.Players) != 0 {
		t.Errorf("Room players count after removal = %v, want 0", len(room.Players))
	}

	// Test removing non-existent player
	err = room.RemovePlayer("non-existent")
	if err == nil {
		t.Error("RemovePlayer() with non-existent ID should return error")
	}
}

// TDD: Test game readiness based on player count
func TestGameReadiness(t *testing.T) {
	room := CreateRoom()

	// Test not ready with 0-2 players
	for i := 0; i < 2; i++ {
		if room.IsReadyToStart() {
			t.Errorf("Room with %d players should not be ready to start", i)
		}

		player := Player{ID: fmt.Sprintf("player-%d", i), Name: fmt.Sprintf("Player%d", i)}
		room.AddPlayer(player)
	}

	// Test not ready with 2 players (need at least 3)
	if room.IsReadyToStart() {
		t.Error("Room with 2 players should not be ready to start")
	}

	// Add third player
	player3 := Player{ID: "player-3", Name: "Player3"}
	room.AddPlayer(player3)

	// Test ready with 3 players
	if !room.IsReadyToStart() {
		t.Error("Room with 3 players should be ready to start")
	}
}

// TDD: Test deck size calculation based on player count
func TestDeckSizeCalculation(t *testing.T) {
	testCases := []struct {
		playerCount       int
		expectedDeckSize  int
		cardsPerInfluence int
	}{
		{3, 15, 3}, // 3-6 players: 3 of each (5 influences)
		{6, 15, 3},
		{7, 20, 4}, // 7-8 players: 4 of each
		{8, 20, 4},
		{9, 25, 5}, // 9-10 players: 5 of each
		{10, 25, 5},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_players", tc.playerCount), func(t *testing.T) {
			deckSize := CalculateDeckSize(tc.playerCount)
			if deckSize != tc.expectedDeckSize {
				t.Errorf("CalculateDeckSize(%d) = %v, want %v",
					tc.playerCount, deckSize, tc.expectedDeckSize)
			}

			cardsPerInfluence := CalculateCardsPerInfluence(tc.playerCount)
			if cardsPerInfluence != tc.cardsPerInfluence {
				t.Errorf("CalculateCardsPerInfluence(%d) = %v, want %v",
					tc.playerCount, cardsPerInfluence, tc.cardsPerInfluence)
			}
		})
	}
}
