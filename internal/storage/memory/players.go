package memory

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/jtbonhomme/gameserver-websocket/internal/players"
)

// ListPlayers returns all registered players (with ID anonymized).
func (m *Memory) ListPlayers() []*players.Player {
	players := []*players.Player{}
	for _, value := range m.players {
		players = append(players, value)
	}
	return players
}

// RegisterPlayer records a player with the given name.
func (m *Memory) RegisterPlayer(id, name string) (*players.Player, error) {
	// if provided id matches a registered player
	if id != uuid.Nil.String() {
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

	player := &players.Player{
		ID:   playerID,
		Name: name,
	}
	m.players[playerID.String()] = player

	return player, nil
}

// UnregisterPlayer removes the player with a given ID.
func (m *Memory) UnregisterPlayer(id string) error {
	_, ok := m.players[id]
	if !ok {
		return fmt.Errorf("unknown player ID: %s", id)
	}
	delete(m.players, id)
	return nil
}
