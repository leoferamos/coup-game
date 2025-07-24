package game

// Card represents a Coup character card with its unique abilities
type Card int

const (
	// Duke allows Tax (3 coins) and blocks Foreign Aid
	Duke Card = iota
	// Assassin allows Assassination (pay 3 coins to eliminate a card)
	Assassin
	// Ambassador allows Exchange (draw 2 cards, keep 2) and blocks Steal
	Ambassador
	// Captain allows Steal (take 2 coins from another player) and blocks Steal
	Captain
	// Contessa blocks Assassination
	Contessa
)

// String returns the string representation of a card
func (c Card) String() string {
	switch c {
	case Duke:
		return "Duke"
	case Assassin:
		return "Assassin"
	case Ambassador:
		return "Ambassador"
	case Captain:
		return "Captain"
	case Contessa:
		return "Contessa"
	default:
		return "Unknown"
	}
}

// GetAllCards returns all available cards in the deck (3 of each type)
func GetAllCards() []Card {
	var deck []Card
	cards := []Card{Duke, Assassin, Ambassador, Captain, Contessa}

	// Add 3 of each card type to the deck
	for _, card := range cards {
		for i := 0; i < 3; i++ {
			deck = append(deck, card)
		}
	}

	return deck
}

// CanPerformAction checks if a card can perform a specific action
func (c Card) CanPerformAction(action ActionType) bool {
	switch action {
	case Tax:
		return c == Duke
	case Assassinate:
		return c == Assassin
	case Exchange:
		return c == Ambassador
	case Steal:
		return c == Captain
	default:
		return false
	}
}

// CanBlock checks if a card can block a specific action
func (c Card) CanBlock(action ActionType) bool {
	switch action {
	case ForeignAid:
		return c == Duke
	case Assassinate:
		return c == Contessa
	case Steal:
		return c == Ambassador || c == Captain
	default:
		return false
	}
}
