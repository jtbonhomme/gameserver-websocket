package players

import (
	"github.com/google/uuid"
)

// Player represents a game player.
type Player struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Score int       `json:"score"`
}

func New(name string) *Player {
	return &Player{
		ID:   uuid.New(),
		Name: name,
	}
}
