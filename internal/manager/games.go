package manager

import (
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/games"
)

// ListGames returns all games.
func (m *Manager) ListGames(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	games := m.store.ListGames()

	b, err = json.Marshal(games)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", games, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = string(b)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

type PayloadGame struct {
	MinPlayers int `json:"minplayers"`
	MaxPlayers int `json:"maxplayers"`
}

// CreateGame instantiates a new game.
func (m *Manager) CreateGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	var game PayloadGame
	err = json.Unmarshal(data, &game)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	createdGame, err := m.store.CreateGame(game.MinPlayers, game.MaxPlayers)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to create game %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	b, err = json.Marshal(createdGame)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = string(b)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)

	_, err = m.node.Publish(ServerPublishChannel, []byte(`{"message": "gameCreated", "data": "`+createdGame.ID().String()+`"}`))
	if err != nil {
		m.log.Error().Msgf("manager publication error: %s", err.Error())
	}
}

// StartGame starts the game with a given ID.
func (m *Manager) StartGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var game games.Game
	err = json.Unmarshal(data, &game)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	err = m.store.StartGame(game.ID().String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to start game %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = EmptyJSON
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

// StopGame stops the game with a given ID.
func (m *Manager) StopGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var game games.Game
	err = json.Unmarshal(data, &game)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	err = m.store.StopGame(game.ID().String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to stop game %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = EmptyJSON
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

// IsGameStarted returns true is game with given ID is started.
func (m *Manager) IsGameStarted(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var started bool
	var err error

	var game games.Game
	err = json.Unmarshal(data, &game)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	started, err = m.store.IsGameStarted(game.ID().String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to stop game %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %t}`, status, started))}, nil)
}

type JoinGameData struct {
	IDGame   uuid.UUID `json:"idGame"`
	IDPlayer uuid.UUID `json:"idPlayer"`
}

// JoinGame adds a player to a game.
func (m *Manager) JoinGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var joinData JoinGameData
	err = json.Unmarshal(data, &joinData)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	err = m.store.JoinGame(joinData.IDGame.String(), joinData.IDPlayer.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to make player %s join game %s: %s", joinData.IDPlayer.String(), joinData.IDGame.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = EmptyJSON
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}
