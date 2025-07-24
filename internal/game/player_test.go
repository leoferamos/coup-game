package game

import (
	"testing"
)

// TDD: Test player creation
func TestNewPlayer(t *testing.T) {
	player := NewPlayer("test-id", "Test Player")

	if player.ID != "test-id" {
		t.Errorf("NewPlayer() ID = %v, want test-id", player.ID)
	}

	if player.Name != "Test Player" {
		t.Errorf("NewPlayer() Name = %v, want Test Player", player.Name)
	}

	if player.Coins != 2 {
		t.Errorf("NewPlayer() Coins = %v, want 2", player.Coins)
	}

	if !player.IsAlive {
		t.Error("NewPlayer() should be alive")
	}

	if !player.IsActive {
		t.Error("NewPlayer() should be active")
	}

	if len(player.Cards) != 0 {
		t.Errorf("NewPlayer() Cards length = %v, want 0", len(player.Cards))
	}
}

// TDD: Test adding cards to player
func TestPlayer_AddCard(t *testing.T) {
	player := NewPlayer("test", "Test")

	// Should be able to add first card
	err := player.AddCard(Duke)
	if err != nil {
		t.Errorf("AddCard() error = %v, want nil", err)
	}

	if len(player.Cards) != 1 {
		t.Errorf("Cards length = %v, want 1", len(player.Cards))
	}

	// Should be able to add second card
	err = player.AddCard(Assassin)
	if err != nil {
		t.Errorf("AddCard() error = %v, want nil", err)
	}

	if len(player.Cards) != 2 {
		t.Errorf("Cards length = %v, want 2", len(player.Cards))
	}

	// Should not be able to add third card
	err = player.AddCard(Captain)
	if err == nil {
		t.Error("AddCard() should return error when adding third card")
	}
}

// TDD: Test removing cards from player
func TestPlayer_RemoveCard(t *testing.T) {
	player := NewPlayer("test", "Test")
	player.AddCard(Duke)
	player.AddCard(Assassin)

	// Should be able to remove existing card
	err := player.RemoveCard(Duke)
	if err != nil {
		t.Errorf("RemoveCard() error = %v, want nil", err)
	}

	if len(player.Cards) != 1 {
		t.Errorf("Cards length = %v, want 1", len(player.Cards))
	}

	if player.HasCard(Duke) {
		t.Error("Player should not have Duke after removal")
	}

	// Should not be able to remove non-existing card
	err = player.RemoveCard(Captain)
	if err == nil {
		t.Error("RemoveCard() should return error for non-existing card")
	}

	// Remove last card should make player dead
	err = player.RemoveCard(Assassin)
	if err != nil {
		t.Errorf("RemoveCard() error = %v, want nil", err)
	}

	if player.IsAlive {
		t.Error("Player should be dead after losing all cards")
	}
}

// TDD: Test checking if player has card
func TestPlayer_HasCard(t *testing.T) {
	player := NewPlayer("test", "Test")
	player.AddCard(Duke)

	if !player.HasCard(Duke) {
		t.Error("HasCard() should return true for Duke")
	}

	if player.HasCard(Assassin) {
		t.Error("HasCard() should return false for Assassin")
	}
}

// TDD: Test coin management
func TestPlayer_CoinManagement(t *testing.T) {
	player := NewPlayer("test", "Test")

	// Test adding coins
	player.AddCoins(3)
	if player.Coins != 5 {
		t.Errorf("Coins = %v, want 5", player.Coins)
	}

	// Test removing coins
	err := player.RemoveCoins(2)
	if err != nil {
		t.Errorf("RemoveCoins() error = %v, want nil", err)
	}

	if player.Coins != 3 {
		t.Errorf("Coins = %v, want 3", player.Coins)
	}

	// Test removing more coins than available
	err = player.RemoveCoins(5)
	if err == nil {
		t.Error("RemoveCoins() should return error for insufficient coins")
	}

	// Test negative coins protection
	player.AddCoins(-10)
	if player.Coins < 0 {
		t.Errorf("Coins = %v, should not be negative", player.Coins)
	}
}

// TDD: Test affordability checks
func TestPlayer_CanAfford(t *testing.T) {
	player := NewPlayer("test", "Test")
	player.Coins = 5

	if !player.CanAfford(Assassinate) { // Costs 3
		t.Error("CanAfford() should return true for Assassinate with 5 coins")
	}

	if player.CanAfford(Coup) { // Costs 7
		t.Error("CanAfford() should return false for Coup with 5 coins")
	}
}

// TDD: Test must coup rule
func TestPlayer_MustCoup(t *testing.T) {
	player := NewPlayer("test", "Test")
	player.Coins = 9

	if player.MustCoup() {
		t.Error("MustCoup() should return false with 9 coins")
	}

	player.Coins = 10
	if !player.MustCoup() {
		t.Error("MustCoup() should return true with 10 coins")
	}
}

// TDD: Test public vs private info
func TestPlayer_GetInfo(t *testing.T) {
	player := NewPlayer("test", "Test Player")
	player.AddCard(Duke)
	player.AddCard(Assassin)

	// Test public info (should not contain cards)
	publicInfo := player.GetPublicInfo()
	if _, hasCards := publicInfo["cards"]; hasCards {
		t.Error("GetPublicInfo() should not contain cards")
	}

	if publicInfo["card_count"] != 2 {
		t.Errorf("GetPublicInfo() card_count = %v, want 2", publicInfo["card_count"])
	}

	// Test private info (should contain cards)
	privateInfo := player.GetPrivateInfo()
	if cards, hasCards := privateInfo["cards"]; !hasCards {
		t.Error("GetPrivateInfo() should contain cards")
	} else if cardSlice, ok := cards.([]string); !ok {
		t.Error("GetPrivateInfo() cards should be []string")
	} else if len(cardSlice) != 2 {
		t.Errorf("GetPrivateInfo() cards length = %v, want 2", len(cardSlice))
	}
}
