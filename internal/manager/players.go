package manager

import (
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
)

// ListPlayers returns the list of all players.
func (m *Manager) ListPlayers(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	players := m.store.ListPlayers()

	for i, _ := range players {
		players[i].ID = uuid.Nil
	}

	b, err = json.Marshal(players)
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

// RegisterPlayer handles new player registration.
func (m *Manager) RegisterPlayer(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	var player players.Player
	err = json.Unmarshal(data, &player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	registeredPlayer, err := m.store.RegisterPlayer(player.ID.String(), player.Name)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to register player %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	b, err = json.Marshal(registeredPlayer)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	_, err = m.node.Publish(ServerPublishChannel, []byte(`{"message": "new player registered", "data": "`+registeredPlayer.Name+`"}`))
	if err != nil {
		m.log.Error().Msgf("manager publication error: %s", err.Error())
	}

	status = OK
	msg = string(b)
	m.log.Debug().Msgf("[rpc] (player) registered: %s", msg)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}

// UnregisterPlayer removes a player from registry.
func (m *Manager) UnregisterPlayer(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var player players.Player
	err = json.Unmarshal(data, &player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	err = m.store.UnregisterPlayer(player.ID.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unregister player %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = "{}"
	m.log.Debug().Msgf("[rpc] player unregistered %s", player.ID)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}
