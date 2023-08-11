package manager

import (
	"github.com/centrifugal/centrifuge"
)

const (
	Register string = "register"
)

func (m *Manager) HandleRPC(e centrifuge.RPCEvent, c centrifuge.RPCCallback) {
	m.log.Info().Msgf("client RPC: %s %s", e.Method, string(e.Data))
	switch e.Method {
	case Register:
		m.Register(e.Data, c)
	default:
		m.log.Error().Msgf("unsupported method %s", e.Method)
	}
	c(centrifuge.RPCReply{Data: []byte(`{"reply": "ok to ` + e.Method + `"}`)}, nil)
}

func (m *Manager) Register(data []byte, c centrifuge.RPCCallback) {
	c(centrifuge.RPCReply{Data: []byte(`{"reply": "client registered"}`)}, nil)
}
