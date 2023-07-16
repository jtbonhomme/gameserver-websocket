package manager

import (
	"fmt"
)

// Player represents a game player.
// todo: ID should be a UUID
type Player struct {
	ID    int
	Name  string
	Score int
}

// RegisterPlayer registers a player with the given name.
func (m *Manager) RegisterPlayer(name string) (*Player, error) {
	result, err := m.db.Exec("INSERT INTO players (name) VALUES (?)", name)
	if err != nil {
		return nil, fmt.Errorf("failed to register player: %v", err)
	}

	playerID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get player ID: %v", err)
	}

	player := &Player{
		ID:   int(playerID),
		Name: name,
	}

	return player, nil
}

// UnregisterPlayer unregisters a player with the given ID.
func (m *Manager) UnregisterPlayer(playerID int) error {
	_, err := m.db.Exec("DELETE FROM players WHERE id = ?", playerID)
	if err != nil {
		return fmt.Errorf("failed to unregister player: %v", err)
	}

	return nil
}
