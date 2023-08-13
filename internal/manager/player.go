package manager

import (
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
)

const (
	OK string = "ok"
	KO string = "ko"
)

// Player represents a game player.
// todo: ID should be a UUID
type Player struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Score int       `json:"score"`
}

// ListAll returns the list of all players.
func (m *Manager) ListAll(data []byte, c centrifuge.RPCCallback) {
	players := m.store.ListAll()
	m.log.Debug().Msgf("[rpc] (player) list all - players %v", players)

	result := fmt.Sprintf(`{"status": %q, "players": [`, OK)
	for i, player := range players {
		result += fmt.Sprintf(`{"name": %q, "score": %d}`, player.Name, player.Score)
		if i < len(players)-1 {
			result += ","
		}
	}
	result += `]}`
	c(centrifuge.RPCReply{Data: []byte(result)}, nil)
}

// Register handles new player registration.
func (m *Manager) Register(data []byte, c centrifuge.RPCCallback) {
	var status, msg string

	m.log.Debug().Msgf("[rpc] (player) register name %s", string(data))
	id := uuid.New()
	var player Player
	err := json.Unmarshal(data, &player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	player.ID = id
	b, err := json.Marshal(player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	_, err = m.store.Register(id, player.Name)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to register player %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = string(b)
	m.log.Debug().Msgf("[rpc] (player) registered: %s", msg)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}

// Unregister removes a player from registry.
func (m *Manager) Unregister(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	m.log.Debug().Msgf("[rpc] (player) unregister - player id %s", string(data))

	var player Player
	err := json.Unmarshal(data, &player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	err = m.store.Unregister(player.ID)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unregister player %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = fmt.Sprintf("unregistered %s (%s)", player.Name, player.ID)
	m.log.Debug().Msgf("[rpc] (player) unregistered: %s", msg)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}
