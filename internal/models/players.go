package models

import (
	"github.com/google/uuid"
)

// Player represents a game player.
// todo: ID should be a UUID
type Player struct {
	ID    uuid.UUID
	Name  string
	Score int
}
