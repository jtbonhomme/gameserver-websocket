package manager

import (
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

const (
	OK string = "ok"
	KO string = "ko"
)

// ListAll returns the list of all players.
func (m *Manager) ListAll(data []byte, c centrifuge.RPCCallback) {
	var status, msg string

	players := m.store.ListAll()

	for i, _ := range players {
		players[i].ID = uuid.Nil
	}

	b, err := json.Marshal(players)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", players, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = string(b)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}

// Register handles new player registration.
func (m *Manager) Register(data []byte, c centrifuge.RPCCallback) {
	var status, msg string

	var player models.Player
	err := json.Unmarshal(data, &player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	registeredPlayer, err := m.store.Register(player.ID.String(), player.Name)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to register player %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	b, err := json.Marshal(registeredPlayer)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", player, err.Error())
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

	var player models.Player
	err := json.Unmarshal(data, &player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	err = m.store.Unregister(player.ID.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unregister player %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = fmt.Sprintf("unregistered: %s (%s)", player.Name, player.ID.String())
	m.log.Debug().Msgf("[rpc] player %s", msg)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}
