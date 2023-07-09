package manager

import (
	"encoding/json"
	"fmt"

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

func (m *Manager) Start() error {
	m.log.Info().Msg("Manager starting ...")

	// connect to server websocket (todo: encapsulate in server package .Dial public function)
	origin := "http://localhost/"
	url := "ws://localhost:12345/connect"
	c, err := websocket.Dial(url, "", origin)
	if err != nil {
		return fmt.Errorf("error connecting to websocket: %w", err)
	}
	m.conn = c

	// subscribe to general topic
	register := pubsub.Message{
		Action: pubsub.SUBSCRIBE,
		Topic:  "com.jtbonhomme.pubsub.general",
	}

	msgR, err := json.Marshal(register)
	if err != nil {
		return fmt.Errorf("error marshaling subscribe message: %w", err)
	}

	_, err = m.conn.Write(msgR)
	if err != nil {
		return fmt.Errorf("error writing subscription message to websocket: %w", err)
	}
	m.started = true

	return nil
}

func (m *Manager) Shutdown() {
	m.log.Info().Msg("Manager shuting down ...")
	m.started = false
}
