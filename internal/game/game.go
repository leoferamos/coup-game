package game

import (
	"fmt"
	"time"
)

// GameState represents the current state of the game
type GameState int

const (
	Waiting GameState = iota
	Starting
	Playing
	Finished
)

// String returns the string representation of game state
func (gs GameState) String() string {
	switch gs {
	case Waiting:
		return "Waiting"
	case Starting:
		return "Starting"
	case Playing:
		return "Playing"
	case Finished:
		return "Finished"
	default:
		return "Unknown"
	}
}

// Game represents a Coup game instance
type Game struct {
	ID            string             `json:"id"`
	State         GameState          `json:"state"`
	Players       map[string]*Player `json:"players"`
	PlayerOrder   []string           `json:"player_order"`
	CurrentPlayer int                `json:"current_player"`
	Deck          []Card             `json:"-"`
	DiscardPile   []Card             `json:"-"`
	CreatedAt     time.Time          `json:"created_at"`
	StartedAt     *time.Time         `json:"started_at,omitempty"`
	FinishedAt    *time.Time         `json:"finished_at,omitempty"`
	Winner        *Player            `json:"winner,omitempty"`
	MinPlayers    int                `json:"min_players"`
	MaxPlayers    int                `json:"max_players"`
}

// NewGame creates a new Coup game instance
func NewGame(id string) *Game {
	return &Game{
		ID:          id,
		State:       Waiting,
		Players:     make(map[string]*Player),
		PlayerOrder: make([]string, 0),
		Deck:        GetAllCards(),
		DiscardPile: make([]Card, 0),
		CreatedAt:   time.Now(),
		MinPlayers:  3,
		MaxPlayers:  6,
	}
}

// AddPlayer adds a player to the game
func (g *Game) AddPlayer(player *Player) error {
	if g.State != Waiting {
		return fmt.Errorf("cannot add players: game is not in waiting state")
	}

	if len(g.Players) >= g.MaxPlayers {
		return fmt.Errorf("game is full: maximum %d players", g.MaxPlayers)
	}

	if _, exists := g.Players[player.ID]; exists {
		return fmt.Errorf("player with ID %s already exists", player.ID)
	}

	g.Players[player.ID] = player
	g.PlayerOrder = append(g.PlayerOrder, player.ID)

	return nil
}

// RemovePlayer removes a player from the game
func (g *Game) RemovePlayer(playerID string) error {
	player, exists := g.Players[playerID]
	if !exists {
		return fmt.Errorf("player with ID %s not found", playerID)
	}

	// If game is playing, mark as inactive instead of removing
	if g.State == Playing {
		player.IsActive = false
		return nil
	}

	// Remove from players map
	delete(g.Players, playerID)

	// Remove from player order
	for i, id := range g.PlayerOrder {
		if id == playerID {
			g.PlayerOrder = append(g.PlayerOrder[:i], g.PlayerOrder[i+1:]...)
			break
		}
	}

	return nil
}

// CanStart checks if the game can be started
func (g *Game) CanStart() bool {
	return g.State == Waiting && len(g.Players) >= g.MinPlayers
}

// StartGame starts the game by shuffling deck and dealing cards
func (g *Game) StartGame() error {
	if !g.CanStart() {
		return fmt.Errorf("cannot start game: need at least %d players", g.MinPlayers)
	}

	g.State = Starting

	// Shuffle the deck
	ShuffleCards(g.Deck)

	// Deal 2 cards to each player
	for _, playerID := range g.PlayerOrder {
		player := g.Players[playerID]
		if err := DealCards(player, &g.Deck); err != nil {
			return fmt.Errorf("failed to deal cards to player %s: %v", playerID, err)
		}
	}

	// Set first player
	g.CurrentPlayer = 0
	g.State = Playing

	now := time.Now()
	g.StartedAt = &now

	return nil
}

// GetCurrentPlayer returns the current player whose turn it is
func (g *Game) GetCurrentPlayer() *Player {
	if g.State != Playing || len(g.PlayerOrder) == 0 {
		return nil
	}

	playerID := g.PlayerOrder[g.CurrentPlayer]
	return g.Players[playerID]
}

// NextTurn advances to the next player's turn
func (g *Game) NextTurn() {
	if g.State != Playing {
		return
	}

	// Find next alive player
	for i := 0; i < len(g.PlayerOrder); i++ {
		g.CurrentPlayer = (g.CurrentPlayer + 1) % len(g.PlayerOrder)
		playerID := g.PlayerOrder[g.CurrentPlayer]
		if g.Players[playerID].IsAlive {
			break
		}
	}

	// Check if game should end
	g.CheckGameEnd()
}

// CheckGameEnd checks if the game should end and sets winner
func (g *Game) CheckGameEnd() {
	if g.State != Playing {
		return
	}

	alivePlayers := g.GetAlivePlayers()

	if len(alivePlayers) <= 1 {
		g.State = Finished
		now := time.Now()
		g.FinishedAt = &now

		if len(alivePlayers) == 1 {
			g.Winner = alivePlayers[0]
		}
	}
}

// GetAlivePlayers returns all players still in the game
func (g *Game) GetAlivePlayers() []*Player {
	var alive []*Player
	for _, player := range g.Players {
		if player.IsAlive {
			alive = append(alive, player)
		}
	}
	return alive
}

// GetGameState returns the current game state for broadcast
func (g *Game) GetGameState() map[string]interface{} {
	players := make([]map[string]interface{}, 0, len(g.PlayerOrder))
	for _, playerID := range g.PlayerOrder {
		players = append(players, g.Players[playerID].GetPublicInfo())
	}

	state := map[string]interface{}{
		"id":             g.ID,
		"state":          g.State.String(),
		"players":        players,
		"current_player": "",
		"deck_size":      len(g.Deck),
	}

	if currentPlayer := g.GetCurrentPlayer(); currentPlayer != nil {
		state["current_player"] = currentPlayer.ID
	}

	if g.Winner != nil {
		state["winner"] = g.Winner.GetPublicInfo()
	}

	return state
}

// GetPlayerGameState returns game state from a specific player's perspective
func (g *Game) GetPlayerGameState(playerID string) map[string]interface{} {
	state := g.GetGameState()

	// Add private information for the requesting player
	if player, exists := g.Players[playerID]; exists {
		state["your_info"] = player.GetPrivateInfo()
		state["your_turn"] = g.GetCurrentPlayer() != nil && g.GetCurrentPlayer().ID == playerID
	}

	return state
}
