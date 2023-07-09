package manager

import (
	"github.com/jtbonhomme/gameserver-websocket/internal/game"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/jtbonhomme/pubsub"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
)

type Manager struct {
	log     *zerolog.Logger
	conn    *websocket.Conn // Client websocket connection.
	games   []*game.Game
	player  []players.Player
	started bool
}

func New(l *zerolog.Logger) *Manager {
	return &Manager{
		log:   l,
		games: []*game.Game{},
	}
}

func (m *Manager) Start() {
	go func() {
		m.log.Info().Msg("Manager starting ...")

		c := pubsub.Client{}
		// connect to server websocket (todo: encapsulate in server package .Dial public function)
		origin := "http://localhost/"
		url := "ws://localhost:12345/connect"
		c.Dial(url, origin)
		c.Register()

		m.started = true
	}()
}

func (m *Manager) Shutdown() {
	// todo stop the goroutine
	m.log.Info().Msg("Manager shuting down ...")
	m.started = false
}
