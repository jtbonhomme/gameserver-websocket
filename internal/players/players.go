package players

import (
	"github.com/google/uuid"
)

// Player represents a game player.
type Player struct {
	ID    uuid.UUID
	Name  string
	Score int
}

func New(name string) *Player {
	return &Player{
		ID:   uuid.New(),
		Name: name,
	}
}
