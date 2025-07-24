package game

import (
	"testing"
)

// TDD: Test ActionType string representation
func TestActionType_String(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected string
	}{
		{Income, "Income"},
		{Coup, "Coup"},
		{ForeignAid, "Foreign Aid"},
		{Tax, "Tax"},
		{Assassinate, "Assassinate"},
		{Exchange, "Exchange"},
		{Steal, "Steal"},
		{ActionType(99), "Unknown"},
	}

	for _, test := range tests {
		if result := test.action.String(); result != test.expected {
			t.Errorf("ActionType.String() = %v, want %v", result, test.expected)
		}
	}
}

// TDD: Test character action identification
func TestActionType_IsCharacterAction(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected bool
	}{
		{Income, false},
		{Coup, false},
		{ForeignAid, false},
		{Tax, true},
		{Assassinate, true},
		{Exchange, true},
		{Steal, true},
	}

	for _, test := range tests {
		if result := test.action.IsCharacterAction(); result != test.expected {
			t.Errorf("ActionType %v.IsCharacterAction() = %v, want %v",
				test.action, result, test.expected)
		}
	}
}

// TDD: Test blockable actions
func TestActionType_CanBeBlocked(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected bool
	}{
		{Income, false},
		{Coup, false},
		{ForeignAid, true},
		{Tax, false},
		{Assassinate, true},
		{Exchange, false},
		{Steal, true},
	}

	for _, test := range tests {
		if result := test.action.CanBeBlocked(); result != test.expected {
			t.Errorf("ActionType %v.CanBeBlocked() = %v, want %v",
				test.action, result, test.expected)
		}
	}
}

// TDD: Test required cards for character actions
func TestActionType_RequiredCard(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected Card
	}{
		{Tax, Duke},
		{Assassinate, Assassin},
		{Exchange, Ambassador},
		{Steal, Captain},
		{Income, -1}, // Invalid for non-character actions
	}

	for _, test := range tests {
		if result := test.action.RequiredCard(); result != test.expected {
			t.Errorf("ActionType %v.RequiredCard() = %v, want %v",
				test.action, result, test.expected)
		}
	}
}

// TDD: Test action costs
func TestActionType_GetCost(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected int
	}{
		{Income, 0},
		{ForeignAid, 0},
		{Tax, 0},
		{Coup, 7},
		{Assassinate, 3},
		{Exchange, 0},
		{Steal, 0},
	}

	for _, test := range tests {
		if result := test.action.GetCost(); result != test.expected {
			t.Errorf("ActionType %v.GetCost() = %v, want %v",
				test.action, result, test.expected)
		}
	}
}

// TDD: Test action rewards
func TestActionType_GetReward(t *testing.T) {
	tests := []struct {
		action   ActionType
		expected int
	}{
		{Income, 1},
		{ForeignAid, 2},
		{Tax, 3},
		{Coup, 0},
		{Assassinate, 0},
		{Exchange, 0},
		{Steal, 2},
	}

	for _, test := range tests {
		if result := test.action.GetReward(); result != test.expected {
			t.Errorf("ActionType %v.GetReward() = %v, want %v",
				test.action, result, test.expected)
		}
	}
}
