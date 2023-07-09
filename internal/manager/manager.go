package manager

import (
	"github.com/jtbonhomme/test-gameserver-websocket/internal/game"
	"github.com/jtbonhomme/test-gameserver-websocket/internal/players"
	"github.com/rs/zerolog"
)

type Manager struct {
	log     *zerolog.Logger
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
	m.log.Info().Msg("Manager starting ...")
	m.started = true
}

func (m *Manager) Shutdown() {
	m.log.Info().Msg("Manager shuting down ...")
	m.started = false
}
