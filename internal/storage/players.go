package storage

import (
	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

// Players defines the interface for players storage.
type Players interface {
	ListAll() []*models.Player                          // ListAll returns all registered players.
	Register(uuid.UUID, string) (*models.Player, error) // Register records a player with the given name.
	Unregister(uuid.UUID) error                         // Unregister removes the player with a given ID.
}
