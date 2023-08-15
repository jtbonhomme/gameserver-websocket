package storage

import (
	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

// Players defines the interface for players storage.
type Players interface {
	ListPlayers() []*models.Player                         // ListPlayers returns all registered players.
	RegisterPlayer(string, string) (*models.Player, error) // Register records a player with the given name.
	UnregisterPlayer(string) error                         // Unregister removes the player with a given ID.
}
