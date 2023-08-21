package players

import (
	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/skyjo"
)

// Player represents a game player.
type Player struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Score int       `json:"score"`
	deck  *skyjo.Deck
}

func New(name string) *Player {
	return &Player{
		ID:   uuid.New(),
		Name: name,
	}
}

func (p *Player) ResetDeck() {
	p.deck.Reset()
}
