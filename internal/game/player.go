package game

import (
	"fmt"
	"math/rand"
	"time"
)

// Player represents a player in the Coup game
type Player struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Coins    int    `json:"coins"`
	Cards    []Card `json:"-"`
	IsAlive  bool   `json:"is_alive"`
	IsActive bool   `json:"is_active"`
}

// NewPlayer creates a new player with starting conditions
func NewPlayer(id, name string) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Coins:    2,
		Cards:    make([]Card, 0, 2),
		IsAlive:  true,
		IsActive: true,
	}
}

// AddCard adds a card to the player's hand
func (p *Player) AddCard(card Card) error {
	if len(p.Cards) >= 2 {
		return fmt.Errorf("player already has maximum cards")
	}
	p.Cards = append(p.Cards, card)
	return nil
}

// RemoveCard removes a specific card from the player's hand
func (p *Player) RemoveCard(card Card) error {
	for i, c := range p.Cards {
		if c == card {
			// Remove card from slice
			p.Cards = append(p.Cards[:i], p.Cards[i+1:]...)

			// Check if player is eliminated
			if len(p.Cards) == 0 {
				p.IsAlive = false
			}

			return nil
		}
	}
	return fmt.Errorf("player does not have card: %s", card.String())
}

// HasCard checks if the player has a specific card
func (p *Player) HasCard(card Card) bool {
	for _, c := range p.Cards {
		if c == card {
			return true
		}
	}
	return false
}

// CanAfford checks if the player can afford an action
func (p *Player) CanAfford(action ActionType) bool {
	return p.Coins >= action.GetCost()
}

// AddCoins adds coins to the player's balance
func (p *Player) AddCoins(amount int) {
	p.Coins += amount
	if p.Coins < 0 {
		p.Coins = 0
	}
}

// RemoveCoins removes coins from the player's balance
func (p *Player) RemoveCoins(amount int) error {
	if p.Coins < amount {
		return fmt.Errorf("insufficient coins: has %d, needs %d", p.Coins, amount)
	}
	p.Coins -= amount
	return nil
}

// MustCoup returns true if player must coup (has 10+ coins)
func (p *Player) MustCoup() bool {
	return p.Coins >= 10
}

// GetPublicInfo returns player information visible to other players
func (p *Player) GetPublicInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"coins":      p.Coins,
		"card_count": len(p.Cards),
		"is_alive":   p.IsAlive,
		"is_active":  p.IsActive,
	}
}

// GetPrivateInfo returns all player information (for the player themselves)
func (p *Player) GetPrivateInfo() map[string]interface{} {
	cardNames := make([]string, len(p.Cards))
	for i, card := range p.Cards {
		cardNames[i] = card.String()
	}

	info := p.GetPublicInfo()
	info["cards"] = cardNames
	return info
}

// ShuffleCards shuffles a slice of cards
func ShuffleCards(cards []Card) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}

// DealCards deals 2 cards to a player from the deck
func DealCards(player *Player, deck *[]Card) error {
	if len(*deck) < 2 {
		return fmt.Errorf("insufficient cards in deck")
	}

	// Deal 2 cards
	for i := 0; i < 2; i++ {
		card := (*deck)[0]
		*deck = (*deck)[1:]

		if err := player.AddCard(card); err != nil {
			return err
		}
	}

	return nil
}
