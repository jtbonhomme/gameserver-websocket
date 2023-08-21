package manager

import (
	"encoding/json"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"

	"github.com/jtbonhomme/gameserver-websocket/internal/games"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

type GameResponse struct {
	ID         uuid.UUID `json:"id"`
	MinPlayers int       `json:"minPlayers"`
	MaxPlayers int       `json:"maxPlayers"`
	Started    bool      `json:"started"`
	TopicName  string    `json:"topicName"`
	Name       string    `json:"name"`
}

// ListGames returns all games.
func (m *Manager) ListGames(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	allGames := []GameResponse{}

	for _, game := range m.games {
		allGames = append(allGames, GameResponse{
			ID:         game.ID,
			MinPlayers: game.MinPlayers,
			MaxPlayers: game.MaxPlayers,
			Started:    game.Started,
			TopicName:  game.TopicName,
			Name:       game.Name,
		})
	}

	b, err = json.Marshal(allGames)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal games data %#v: %s", allGames, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = string(b)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

type GameRequest struct {
	ID         uuid.UUID `json:"id"`
	MinPlayers int       `json:"minPlayers"`
	MaxPlayers int       `json:"maxPlayers"`
}

// CreateGame instantiates a new game.
func (m *Manager) CreateGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var b []byte
	var err error

	var req GameRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	game := games.New(m.log, req.MinPlayers, req.MaxPlayers)

	b, err = json.Marshal(game)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to marshal data %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	m.games[game.ID.String()] = game

	status = OK
	msg = string(b)
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)

	_, err = m.node.Publish(utils.ServerPublishChannel,
		[]byte(`{"type": "creation", "emitter": "manager", "id": "`+game.ID.String()+`", "data": ""}`))
	if err != nil {
		m.log.Error().Msgf("manager publication error: %s", err.Error())
	}
}

// StartGame starts the game with a given ID.
func (m *Manager) StartGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var req GameRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	game, err := m.gameByID(req.ID.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to retrieve game from its ID: %s", err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	err = game.Start()
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to start game %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	m.games[game.ID.String()] = game

	// publication to all clients who subscribed to a channel
	_, err = m.node.Publish(utils.ServerPublishChannel,
		[]byte(`{"type": "start", "emitter": "manager", "id": "`+game.ID.String()+`", "data": ""}`))
	if err != nil {
		m.log.Error().Msgf("manager publication error: %s", err.Error())
	}
	/*
		// message to one client
		for _, playerID := range game.Players() {
			m.playersToClientsMap[playerID].Send([]byte(`{"id": "` + playerID + `", "action": "do something"}`))
		}
	*/
	status = OK
	msg = EmptyJSON
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

// StopGame stops the game with a given ID.
func (m *Manager) StopGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var req GameRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	game, err := m.gameByID(req.ID.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to retrieve game from its ID: %s", err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	err = game.Stop()
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to stop game %v: %s", game, err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	m.games[game.ID.String()] = game

	status = OK
	msg = EmptyJSON
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

type GamePlayerRequest struct {
	IDGame   uuid.UUID `json:"idGame"`
	IDPlayer uuid.UUID `json:"idPlayer"`
}

// JoinGame adds a player to a game.
func (m *Manager) JoinGame(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var req GamePlayerRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	game, err := m.gameByID(req.IDGame.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to retrieve game from its ID %s: %s", req.IDGame.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	player, err := m.playerByID(req.IDPlayer.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to retrieve player from its ID %s: %s", req.IDPlayer.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	err = game.AddPlayer(player)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to make player %s join game %s: %s", req.IDPlayer.String(), req.IDGame.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	m.games[game.ID.String()] = game

	_, err = m.node.Publish(utils.ServerPublishChannel,
		[]byte(`{"type": "join", "emitter": "manager", "id": "`+req.IDGame.String()+`", "data": "`+player.Name+`"}`))
	if err != nil {
		m.log.Error().Msgf("manager publication error: %s", err.Error())
	}

	status = OK
	msg = EmptyJSON
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

// PlayerInit starts the game with a given ID.
func (m *Manager) PlayerInit(data []byte, c centrifuge.RPCCallback) {
	var status, msg string
	var err error

	var initData GamePlayerRequest
	err = json.Unmarshal(data, &initData)
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to unmarshal data %q: %s", string(data), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	game, err := m.gameByID(initData.IDGame.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to retrieve game from its ID %s: %s", initData.IDGame.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	err = game.PlayerInit(initData.IDPlayer.String())
	if err != nil {
		status = KO
		msg = fmt.Sprintf("unable to init player %s: %s", initData.IDPlayer.String(), err.Error())
		c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
		return
	}

	status = OK
	msg = EmptyJSON
	c(centrifuge.RPCReply{Data: []byte(fmt.Sprintf(`{"status": %q, "result": %q}`, status, msg))}, nil)
}

// gameByID returns a game object from its ID.
func (m *Manager) gameByID(id string) (*games.Game, error) {
	// if provided id matches an existing game
	if id != uuid.Nil.String() {
		game, ok := m.games[id]
		if ok {
			return game, nil
		}
	}

	return nil, fmt.Errorf("unknown game id: %s", id)
}
