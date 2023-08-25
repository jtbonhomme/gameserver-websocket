package games

import (
	"fmt"

	"github.com/google/uuid"
)

// playerByID returns a player index from its ID.
func (game *Game) playerByID(id string) (int, error) {
	// if provided id matches a registered player
	if id != uuid.Nil.String() {
		for i, p := range game.players {
			if p.ID.String() == id {
				return i, nil
			}
		}
	}

	return 0, fmt.Errorf("unknown player id: %s", id)
}
