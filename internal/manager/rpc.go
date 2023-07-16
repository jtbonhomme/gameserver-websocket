package manager

import (
	"github.com/centrifugal/centrifuge"
)

func (m *Manager) HandleRPC(e centrifuge.RPCEvent, c centrifuge.RPCCallback) {
	m.log.Info().Msgf("client RPC: %s %s", e.Method, string(e.Data))
	c(centrifuge.RPCReply{Data: []byte(`{"reply": "ok to ` + e.Method + `"}`)}, nil)
}
