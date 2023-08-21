package manager

import (
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"

	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

type PlayerResponse struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

// ListPlayers returns the list of all players.
func (m *Manager) ListPlayers(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	allPlayers := []PlayerResponse{}

	for _, player := range m.players {
		allPlayers = append(allPlayers, PlayerResponse{
			Name:  player.Name,
			Score: player.Score,
		})
	}

	b, err = json.Marshal(allPlayers)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", allPlayers, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = string(b)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}

type PlayerRequest struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// RegisterPlayer handles new player registration.
func (m *Manager) RegisterPlayer(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	var req PlayerRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	player := m.registerPlayer(req.ID.String(), req.Name)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to register player %s / %s: %s", req.Name, req.ID.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	b, err = json.Marshal(player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", player, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	_, err = m.node.Publish(utils.ServerPublishChannel,
		[]byte(`{"type": "registration", "emitter": "manager", "id": "", "data": "`+player.Name+`"}`))
	if err != nil {
		m.log.Error().Msgf("manager publication error: %s", err.Error())
	}

	status = OK
	msg = string(b)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}

// UnregisterPlayer removes a player from registry.
func (m *Manager) UnregisterPlayer(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var req PlayerRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	err = m.unregisterPlayer(req.ID.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unregister player %s: %s", req.ID.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = "{}"
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result":%q}`, status, msg))}, nil)
}

// playerByID returns a player object from its ID.
func (m *Manager) playerByID(id string) (*players.Player, error) {
	// if provided id matches a registered player
	if id != uuid.Nil.String() {
		player, ok := m.players[id]
		if ok {
			return player, nil
		}
	}

	return nil, fmt.Errorf("unknown player id: %s", id)
}

// registerPlayer records a player with the given name.
func (m *Manager) registerPlayer(id, name string) *players.Player {
	// if provided id matches a registered player
	if id != uuid.Nil.String() {
		p, ok := m.players[id]
		if ok {
			return p
		}
	}

	player := players.New(name)
	m.players[player.ID.String()] = player

	return player
}

// unregisterPlayer removes the player with a given ID.
func (m *Manager) unregisterPlayer(id string) error {
	_, ok := m.players[id]
	if !ok {
		return fmt.Errorf("unknown player ID: %s", id)
	}

	delete(m.players, id)

	return nil
}
