package memory

import (
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

type Memory struct {
	log     *zerolog.Logger
	players map[string]*models.Player
	games   map[string]*models.Game
}

// New creates a new Memory object.
func New(l *zerolog.Logger) *Memory {
	mem := &Memory{
		log:     l,
		players: make(map[string]*models.Player),
		games:   make(map[string]*models.Game),
	}

	return mem
}
