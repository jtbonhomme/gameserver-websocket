package manager

import (
	"fmt"

	"github.com/centrifugal/centrifuge"
)

const (
	RegisterPlayer   string = "registerPlayer"
	UnregisterPlayer string = "unregisterPlayer"
	ListPlayers      string = "listPlayers"
	ListGames        string = "listGames"
	CreateGame       string = "createGame"
	StartGame        string = "startGame"
	StopGame         string = "stopGame"
	IsGameStarted    string = "isGameStarted"
	JoinGame         string = "joinGame"
)

const (
	OK        string = "ok"
	KO        string = "ko"
	EmptyJSON string = "{}"
)

// HandleRPC execute remote procedure call defined by the RPCEvent, then call the provided callback.
func (m *Manager) HandleRPC(e centrifuge.RPCEvent, c centrifuge.RPCCallback) {
	m.log.Info().Msgf("client RPC: %s %s", e.Method, string(e.Data))
	// Players related rpc
	switch e.Method {
	case RegisterPlayer:
		m.RegisterPlayer(e.Data, c)
	case UnregisterPlayer:
		m.UnregisterPlayer(e.Data, c)
	case ListPlayers:
		m.ListPlayers(e.Data, c)
	// Games related rpc
	case ListGames:
		m.ListGames(e.Data, c)
	case CreateGame:
		m.CreateGame(e.Data, c)
	case StartGame:
		m.StartGame(e.Data, c)
	case StopGame:
		m.StopGame(e.Data, c)
	case IsGameStarted:
		m.IsGameStarted(e.Data, c)
	case JoinGame:
		m.JoinGame(e.Data, c)
	// Default
	default:
		msg := fmt.Sprintf("unsupported method %s", e.Method)
		m.log.Error().Msg(msg)
		c(centrifuge.RPCReply{Data: []byte(`{"status": "ko", "reason":"` + msg + `"}`)}, nil)
	}
}
