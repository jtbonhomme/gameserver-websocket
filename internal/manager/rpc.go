package manager

import (
	"github.com/centrifugal/centrifuge"
)

const (
	Register string = "register"
)

// HandleRPC execute remote procedure call defined by the RPCEvent, then call the provided callback.
func (m *Manager) HandleRPC(e centrifuge.RPCEvent, c centrifuge.RPCCallback) {
	m.log.Info().Msgf("client RPC: %s %s", e.Method, string(e.Data))
	switch e.Method {
	case Register:
		m.Register(e.Data, c)
	default:
		m.log.Error().Msgf("unsupported method %s", e.Method)
	}
}
