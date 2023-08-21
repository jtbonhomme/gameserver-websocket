package memory

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/players"
)

type Memory struct {
	log     *zerolog.Logger
	players map[string]*players.Player
}

// New creates a new Memory object.
func New(l *zerolog.Logger) *Memory {
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[memory] %s", i) },
	}
	log := l.Output(output)

	mem := &Memory{
		log:     &log,
		players: make(map[string]*players.Player),
	}

	return mem
}
