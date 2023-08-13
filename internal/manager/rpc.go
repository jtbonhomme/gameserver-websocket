package manager

import (
	"fmt"

	"github.com/centrifugal/centrifuge"
)

const (
	Register   string = "register"
	Unregister string = "unregister"
	ListAll    string = "listAll"
)

// HandleRPC execute remote procedure call defined by the RPCEvent, then call the provided callback.
func (m *Manager) HandleRPC(e centrifuge.RPCEvent, c centrifuge.RPCCallback) {
	m.log.Info().Msgf("client RPC: %s %s", e.Method, string(e.Data))
	switch e.Method {
	case Register:
		m.Register(e.Data, c)
	case Unregister:
		m.Unregister(e.Data, c)
	case ListAll:
		m.ListAll(e.Data, c)
	default:
		msg := fmt.Sprintf("unsupported method %s", e.Method)
		m.log.Error().Msg(msg)
		c(centrifuge.RPCReply{Data: []byte(`{"status": "ko", "reason":"` + msg + `"}`)}, nil)
	}
}
