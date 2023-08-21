package players

import "github.com/google/uuid"

// ContainsID checks if a player is already present in a slice.
func ContainsID(players []*Player, id uuid.UUID) bool {
	for _, p := range players {
		if p.ID.String() == id.String() {
			return true
		}
	}

	return false
}
