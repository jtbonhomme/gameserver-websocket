package manager

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/jtbonhomme/gameserver-websocket/internal/game"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/jtbonhomme/pubsub"
	"github.com/rs/zerolog"
	"golang.org/x/net/websocket"
)

const (
	_defaultShutdownTimeout = 3 * time.Second
)

type Manager struct {
	log             *zerolog.Logger
	conn            *websocket.Conn // Client websocket connection.
	games           []*game.Game
	player          []players.Player
	started         bool
	err             chan error
	shutdownTimeout time.Duration
	shutdownChannel chan struct{}
	waitGroup       *sync.WaitGroup
}

func New(l *zerolog.Logger) *Manager {
	return &Manager{
		log:             l,
		games:           []*game.Game{},
		shutdownTimeout: _defaultShutdownTimeout,
		shutdownChannel: make(chan struct{}),
		waitGroup:       &sync.WaitGroup{},
	}
}

func (m *Manager) Start() {

	m.waitGroup.Add(1)
	var err error
	go func(shutdownChannel chan struct{}, wg *sync.WaitGroup) {
		m.log.Info().Msg("Manager starting ...")
		defer wg.Done()
		c := pubsub.Client{}
		// connect to server websocket
		origin := "http://localhost/"
		url := "ws://localhost:12345/connect"
		err = c.Dial(url, origin)
		if err != nil {
			m.err <- err
		}
		err = c.Register("com.jtbonhomme.pubsub.general")
		if err != nil {
			m.err <- err
		}

		m.started = true
		for {
			select {
			case <-shutdownChannel:
				m.log.Info().Msg("Shutdown manager goroutine")
				return
			default:
				runtime.Gosched()
			}
		}

	}(m.shutdownChannel, m.waitGroup)
}

// Error -.
func (m *Manager) Error() <-chan error {
	return m.err
}

func (m *Manager) Shutdown() {
	_, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
	defer cancel()

	m.log.Info().Msg("Manager shuting down ...")
	m.started = false
	close(m.shutdownChannel)
	m.waitGroup.Wait() // wait for all goroutines
}
