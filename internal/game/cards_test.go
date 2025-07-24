package game

import (
	"testing"
)

// TDD: Test Card string representation
func TestCard_String(t *testing.T) {
	tests := []struct {
		card     Card
		expected string
	}{
		{Duke, "Duke"},
		{Assassin, "Assassin"},
		{Ambassador, "Ambassador"},
		{Captain, "Captain"},
		{Contessa, "Contessa"},
		{Card(99), "Unknown"},
	}

	for _, test := range tests {
		if result := test.card.String(); result != test.expected {
			t.Errorf("Card.String() = %v, want %v", result, test.expected)
		}
	}
}

// TDD: Test deck creation
func TestGetAllCards(t *testing.T) {
	deck := GetAllCards()

	// Should have 15 cards total (3 of each type)
	if len(deck) != 15 {
		t.Errorf("GetAllCards() length = %v, want 15", len(deck))
	}

	// Count each card type
	counts := make(map[Card]int)
	for _, card := range deck {
		counts[card]++
	}

	expectedCards := []Card{Duke, Assassin, Ambassador, Captain, Contessa}
	for _, card := range expectedCards {
		if counts[card] != 3 {
			t.Errorf("Card %v count = %v, want 3", card, counts[card])
		}
	}
}

// TDD: Test card abilities
func TestCard_CanPerformAction(t *testing.T) {
	tests := []struct {
		card     Card
		action   ActionType
		expected bool
	}{
		{Duke, Tax, true},
		{Duke, Steal, false},
		{Assassin, Assassinate, true},
		{Assassin, Tax, false},
		{Ambassador, Exchange, true},
		{Captain, Steal, true},
		{Contessa, Tax, false},
	}

	for _, test := range tests {
		if result := test.card.CanPerformAction(test.action); result != test.expected {
			t.Errorf("Card %v.CanPerformAction(%v) = %v, want %v",
				test.card, test.action, result, test.expected)
		}
	}
}

// TDD: Test card blocking abilities
func TestCard_CanBlock(t *testing.T) {
	tests := []struct {
		card     Card
		action   ActionType
		expected bool
	}{
		{Duke, ForeignAid, true},
		{Duke, Steal, false},
		{Contessa, Assassinate, true},
		{Ambassador, Steal, true},
		{Captain, Steal, true},
		{Assassin, ForeignAid, false},
	}

	for _, test := range tests {
		if result := test.card.CanBlock(test.action); result != test.expected {
			t.Errorf("Card %v.CanBlock(%v) = %v, want %v",
				test.card, test.action, result, test.expected)
		}
	}
}
