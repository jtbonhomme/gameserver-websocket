package manager

import (
	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
)

// Player represents a game player.
// todo: ID should be a UUID
type Player struct {
	ID    uuid.UUID
	Name  string
	Score int
}

// ListAll handles all players listing.
func (m *Manager) ListAll(data []byte, c centrifuge.RPCCallback) {
	c(centrifuge.RPCReply{Data: []byte(`{"status": "ok", "id":"` + `"}`)}, nil)
}

// Register handles new player registration.
func (m *Manager) Register(data []byte, c centrifuge.RPCCallback) {
	m.log.Debug().Msgf("[rpc] (player) register - %s", string(data))
	c(centrifuge.RPCReply{Data: []byte(`{"status": "ok", "id":"` + `"}`)}, nil)
}

// Unregister handles new player removal from registry.
func (m *Manager) Unregister(data []byte, c centrifuge.RPCCallback) {
	c(centrifuge.RPCReply{Data: []byte(`{"status": "ok", "id":"` + `"}`)}, nil)
}
