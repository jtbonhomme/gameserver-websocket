package storage

import (
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
)

// Players defines the interface for players storage.
type Players interface {
	ListPlayers() []*players.Player                         // ListPlayers returns all registered players.
	RegisterPlayer(string, string) (*players.Player, error) // Register records a player with the given name.
	UnregisterPlayer(string) error                          // Unregister removes the player with a given ID.
	NameByID(string) (string, error)                        // NameByID returns a player name from its ID.
}
