package sqlite

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

// ListAll returns all registered players.
func (s *SQLite) ListAll() []*models.Player {
	return []*models.Player{}
}

// Register registers a player with the given name.
func (s *SQLite) Register(id uuid.UUID, name string) (*models.Player, error) {
	playerID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to create player ID: %v", err)
	}

	player := &models.Player{
		ID:   playerID,
		Name: name,
	}

	_, err = s.db.Exec("INSERT INTO players (name, 0) VALUES (?)", name)
	if err != nil {
		return nil, fmt.Errorf("failed to register player: %v", err)
	}

	return player, nil
}

// Unregister unregisters a player with the given ID.
func (s *SQLite) Unregister(playerID uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM players WHERE id = ?", playerID.String())
	if err != nil {
		return fmt.Errorf("failed to unregister player: %v", err)
	}

	return nil
}
