package game

import (
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/rs/zerolog"
)

type Game struct {
	log     *zerolog.Logger
	started bool
	players []players.Player
}

func New(l *zerolog.Logger) *Game {
	return &Game{
		log:     l,
		players: []players.Player{},
	}
}

func (g *Game) Start() {
	g.started = true
}

func (g *Game) Stop() {
	g.started = false
}
