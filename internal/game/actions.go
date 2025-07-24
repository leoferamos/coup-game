package game

// ActionType represents the different types of actions a player can take
type ActionType int

const (
	// Basic actions (always available)
	Income ActionType = iota // Take 1 coin
	Coup                     // Pay 7 coins to eliminate any player's card

	// Actions that can be blocked
	ForeignAid // Take 2 coins (can be blocked by Duke)

	// Character-specific actions
	Tax         // Duke: Take 3 coins
	Assassinate // Assassin: Pay 3 coins to eliminate a card (can be blocked by Contessa)
	Exchange    // Ambassador: Draw 2 cards, keep 2, return 2 to deck
	Steal       // Captain: Take 2 coins from another player (can be blocked by Captain/Ambassador)
)

// String returns the string representation of an action
func (a ActionType) String() string {
	switch a {
	case Income:
		return "Income"
	case Coup:
		return "Coup"
	case ForeignAid:
		return "Foreign Aid"
	case Tax:
		return "Tax"
	case Assassinate:
		return "Assassinate"
	case Exchange:
		return "Exchange"
	case Steal:
		return "Steal"
	default:
		return "Unknown"
	}
}

// IsCharacterAction returns true if the action requires a specific character card
func (a ActionType) IsCharacterAction() bool {
	return a == Tax || a == Assassinate || a == Exchange || a == Steal
}

// CanBeBlocked returns true if the action can be blocked by another player
func (a ActionType) CanBeBlocked() bool {
	return a == ForeignAid || a == Assassinate || a == Steal
}

// RequiredCard returns the card needed to perform this character action
func (a ActionType) RequiredCard() Card {
	switch a {
	case Tax:
		return Duke
	case Assassinate:
		return Assassin
	case Exchange:
		return Ambassador
	case Steal:
		return Captain
	default:
		return -1 // Invalid card for non-character actions
	}
}

// GetCost returns the coin cost of the action
func (a ActionType) GetCost() int {
	switch a {
	case Coup:
		return 7
	case Assassinate:
		return 3
	default:
		return 0
	}
}

// GetReward returns the coin reward of the action
func (a ActionType) GetReward() int {
	switch a {
	case Income:
		return 1
	case ForeignAid:
		return 2
	case Tax:
		return 3
	case Steal:
		return 2 // Stolen from another player
	default:
		return 0
	}
}
