package game

import (
	"fmt"
	"testing"
)

// TDD: Test game creation
func TestNewGame(t *testing.T) {
	game := NewGame("test-game")

	if game.ID != "test-game" {
		t.Errorf("NewGame() ID = %v, want test-game", game.ID)
	}

	if game.State != Waiting {
		t.Errorf("NewGame() State = %v, want Waiting", game.State)
	}

	if len(game.Players) != 0 {
		t.Errorf("NewGame() Players length = %v, want 0", len(game.Players))
	}

	if len(game.Deck) != 15 {
		t.Errorf("NewGame() Deck length = %v, want 15", len(game.Deck))
	}

	if game.MinPlayers != 3 {
		t.Errorf("NewGame() MinPlayers = %v, want 3", game.MinPlayers)
	}

	if game.MaxPlayers != 6 {
		t.Errorf("NewGame() MaxPlayers = %v, want 6", game.MaxPlayers)
	}
}

// TDD: Test adding players to game
func TestGame_AddPlayer(t *testing.T) {
	game := NewGame("test")
	player1 := NewPlayer("p1", "Player 1")
	player2 := NewPlayer("p2", "Player 2")

	// Should be able to add first player
	err := game.AddPlayer(player1)
	if err != nil {
		t.Errorf("AddPlayer() error = %v, want nil", err)
	}

	if len(game.Players) != 1 {
		t.Errorf("Players length = %v, want 1", len(game.Players))
	}

	if len(game.PlayerOrder) != 1 {
		t.Errorf("PlayerOrder length = %v, want 1", len(game.PlayerOrder))
	}

	// Should be able to add second player
	err = game.AddPlayer(player2)
	if err != nil {
		t.Errorf("AddPlayer() error = %v, want nil", err)
	}

	// Should not be able to add duplicate player
	err = game.AddPlayer(player1)
	if err == nil {
		t.Error("AddPlayer() should return error for duplicate player")
	}
}

// TDD: Test game cannot add players when not waiting
func TestGame_AddPlayer_NotWaiting(t *testing.T) {
	game := NewGame("test")
	game.State = Playing
	player := NewPlayer("p1", "Player 1")

	err := game.AddPlayer(player)
	if err == nil {
		t.Error("AddPlayer() should return error when game is not waiting")
	}
}

// TDD: Test removing players from game
func TestGame_RemovePlayer(t *testing.T) {
	game := NewGame("test")
	player := NewPlayer("p1", "Player 1")
	game.AddPlayer(player)

	// Should be able to remove existing player in waiting state
	err := game.RemovePlayer("p1")
	if err != nil {
		t.Errorf("RemovePlayer() error = %v, want nil", err)
	}

	if len(game.Players) != 0 {
		t.Errorf("Players length = %v, want 0", len(game.Players))
	}

	// Should return error for non-existing player
	err = game.RemovePlayer("nonexistent")
	if err == nil {
		t.Error("RemovePlayer() should return error for non-existing player")
	}
}

// TDD: Test removing player during game marks as inactive
func TestGame_RemovePlayer_Playing(t *testing.T) {
	game := NewGame("test")
	player := NewPlayer("p1", "Player 1")
	game.AddPlayer(player)
	game.State = Playing

	err := game.RemovePlayer("p1")
	if err != nil {
		t.Errorf("RemovePlayer() error = %v, want nil", err)
	}

	// Player should still exist but be inactive
	if len(game.Players) != 1 {
		t.Errorf("Players length = %v, want 1", len(game.Players))
	}

	if game.Players["p1"].IsActive {
		t.Error("Player should be inactive after removal during game")
	}
}

// TDD: Test game start conditions
func TestGame_CanStart(t *testing.T) {
	game := NewGame("test")

	// Should not be able to start with no players
	if game.CanStart() {
		t.Error("CanStart() should return false with no players")
	}

	// Add minimum players
	for i := 0; i < 3; i++ {
		player := NewPlayer(fmt.Sprintf("p%d", i), fmt.Sprintf("Player %d", i))
		game.AddPlayer(player)
	}

	// Should be able to start with minimum players
	if !game.CanStart() {
		t.Error("CanStart() should return true with 3 players")
	}

	// Should not be able to start if not in waiting state
	game.State = Playing
	if game.CanStart() {
		t.Error("CanStart() should return false when not in waiting state")
	}
}

// TDD: Test starting the game
func TestGame_StartGame(t *testing.T) {
	game := NewGame("test")

	// Add enough players
	for i := 0; i < 3; i++ {
		player := NewPlayer(fmt.Sprintf("p%d", i), fmt.Sprintf("Player %d", i))
		game.AddPlayer(player)
	}

	err := game.StartGame()
	if err != nil {
		t.Errorf("StartGame() error = %v, want nil", err)
	}

	// Check game state changed
	if game.State != Playing {
		t.Errorf("State = %v, want Playing", game.State)
	}

	// Check players have cards
	for _, player := range game.Players {
		if len(player.Cards) != 2 {
			t.Errorf("Player %s has %d cards, want 2", player.ID, len(player.Cards))
		}
	}

	// Check current player is set
	if game.CurrentPlayer != 0 {
		t.Errorf("CurrentPlayer = %v, want 0", game.CurrentPlayer)
	}

	// Check started time is set
	if game.StartedAt == nil {
		t.Error("StartedAt should be set")
	}
}

// TDD: Test getting current player
func TestGame_GetCurrentPlayer(t *testing.T) {
	game := NewGame("test")

	// Should return nil when not playing
	if current := game.GetCurrentPlayer(); current != nil {
		t.Error("GetCurrentPlayer() should return nil when not playing")
	}

	// Add players and start game
	for i := 0; i < 3; i++ {
		player := NewPlayer(fmt.Sprintf("p%d", i), fmt.Sprintf("Player %d", i))
		game.AddPlayer(player)
	}
	game.StartGame()

	// Should return first player
	current := game.GetCurrentPlayer()
	if current == nil {
		t.Error("GetCurrentPlayer() should not return nil during game")
		return
	}

	if current.ID != "p0" {
		t.Errorf("GetCurrentPlayer() ID = %v, want p0", current.ID)
	}
}

// TDD: Test turn advancement
func TestGame_NextTurn(t *testing.T) {
	game := NewGame("test")

	// Add players and start game
	for i := 0; i < 3; i++ {
		player := NewPlayer(fmt.Sprintf("p%d", i), fmt.Sprintf("Player %d", i))
		game.AddPlayer(player)
	}
	game.StartGame()

	// Should advance to next player
	game.NextTurn()
	current := game.GetCurrentPlayer()
	if current == nil {
		t.Error("GetCurrentPlayer() should not return nil after NextTurn()")
		return
	}
	if current.ID != "p1" {
		t.Errorf("GetCurrentPlayer() after NextTurn() = %v, want p1", current.ID)
	}

	// Should wrap around
	game.NextTurn()
	game.NextTurn()
	current = game.GetCurrentPlayer()
	if current == nil {
		t.Error("GetCurrentPlayer() should not return nil after wrap around")
		return
	}
	if current.ID != "p0" {
		t.Errorf("GetCurrentPlayer() after wrap around = %v, want p0", current.ID)
	}
}

// TDD: Test game end detection
func TestGame_CheckGameEnd(t *testing.T) {
	game := NewGame("test")

	// Add players and start game
	for i := 0; i < 3; i++ {
		player := NewPlayer(fmt.Sprintf("p%d", i), fmt.Sprintf("Player %d", i))
		game.AddPlayer(player)
	}
	game.StartGame()

	// Eliminate players until one remains
	game.Players["p0"].IsAlive = false
	game.Players["p1"].IsAlive = false

	game.CheckGameEnd()

	// Game should be finished
	if game.State != Finished {
		t.Errorf("State = %v, want Finished", game.State)
	}

	// Winner should be set
	if game.Winner == nil {
		t.Error("Winner should be set")
		return
	}

	if game.Winner.ID != "p2" {
		t.Errorf("Winner ID = %v, want p2", game.Winner.ID)
	}

	// Finished time should be set
	if game.FinishedAt == nil {
		t.Error("FinishedAt should be set")
	}
}

// TDD: Test getting alive players
func TestGame_GetAlivePlayers(t *testing.T) {
	game := NewGame("test")

	// Add players
	for i := 0; i < 3; i++ {
		player := NewPlayer(fmt.Sprintf("p%d", i), fmt.Sprintf("Player %d", i))
		game.AddPlayer(player)
	}

	// All should be alive initially
	alive := game.GetAlivePlayers()
	if len(alive) != 3 {
		t.Errorf("GetAlivePlayers() length = %v, want 3", len(alive))
	}

	// Eliminate one player
	game.Players["p1"].IsAlive = false
	alive = game.GetAlivePlayers()
	if len(alive) != 2 {
		t.Errorf("GetAlivePlayers() length = %v, want 2", len(alive))
	}
}

// TDD: Test game state serialization
func TestGame_GetGameState(t *testing.T) {
	game := NewGame("test")
	player := NewPlayer("p1", "Player 1")
	game.AddPlayer(player)

	state := game.GetGameState()

	// Check basic fields
	if state["id"] != "test" {
		t.Errorf("Game state id = %v, want test", state["id"])
	}

	if state["state"] != "Waiting" {
		t.Errorf("Game state state = %v, want Waiting", state["state"])
	}

	// Check players array
	if players, ok := state["players"].([]map[string]interface{}); !ok {
		t.Error("Game state players should be []map[string]interface{}")
	} else if len(players) != 1 {
		t.Errorf("Game state players length = %v, want 1", len(players))
	}
}

// TDD: Test player-specific game state
func TestGame_GetPlayerGameState(t *testing.T) {
	game := NewGame("test")
	player := NewPlayer("p1", "Player 1")
	game.AddPlayer(player)

	state := game.GetPlayerGameState("p1")

	// Should contain private player info
	if _, hasYourInfo := state["your_info"]; !hasYourInfo {
		t.Error("Player game state should contain your_info")
	}

	if _, hasYourTurn := state["your_turn"]; !hasYourTurn {
		t.Error("Player game state should contain your_turn")
	}
}
