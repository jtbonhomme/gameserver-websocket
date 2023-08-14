package storage

import (
	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

// Players defines the interface for players storage.
type Players interface {
	ListAll() []*models.Player                       // ListAll returns all registered players.
	Register(string, string) (*models.Player, error) // Register records a player with the given name.
	Unregister(string) error                         // Unregister removes the player with a given ID.
}
