package memory

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

type Memory struct {
	players map[uuid.UUID]*models.Player
}

// New creates a new Memory object.
func New() *Memory {
	mem := &Memory{
		players: make(map[uuid.UUID]*models.Player),
	}

	return mem
}

// ListAll returns all registered players (with ID anonymized).
func (m *Memory) ListAll() []*models.Player {
	players := []*models.Player{}
	for _, value := range m.players {
		value.ID = uuid.Nil // anonymize player's ID.
		players = append(players, value)
	}

	return players
}

// Register records a player with the given name.
func (m *Memory) Register(id uuid.UUID, name string) (*models.Player, error) {
	// if provided id matches a registered player
	if id != uuid.Nil {
		player, ok := m.players[id]
		if ok {
			return player, nil
		}
	}
	// else, create a new register for this player
	playerID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to create player ID: %v", err)
	}

	player := &models.Player{
		ID:   playerID,
		Name: name,
	}
	m.players[playerID] = player

	return player, nil
}

// Unregister removes the player with a given ID.
func (m *Memory) Unregister(id uuid.UUID) error {
	_, ok := m.players[id]
	if !ok {
		return fmt.Errorf("unknown player ID: %s", id.String())
	}
	delete(m.players, id)
	return nil
}
