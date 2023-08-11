package manager

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/jtbonhomme/gameserver-websocket/internal/game"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/jtbonhomme/pubsub"
	"github.com/rs/zerolog"
)

const (
	_defaultShutdownTimeout = 3 * time.Second
)

type Manager struct {
	log             *zerolog.Logger
	games           []*game.Game
	player          []players.Player
	started         bool
	err             chan error
	shutdownTimeout time.Duration
	shutdownChannel chan struct{} // Channel used to stop manager.
	waitGroup       *sync.WaitGroup
	client          *pubsub.Client
}

func New(l *zerolog.Logger) *Manager {
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[manager] %s", i) },
	}
	logger := l.Output(output)
	return &Manager{
		log:             &logger,
		games:           []*game.Game{},
		shutdownTimeout: _defaultShutdownTimeout,
		shutdownChannel: make(chan struct{}),
		err:             make(chan error),
		waitGroup:       &sync.WaitGroup{},
	}
}

func (m *Manager) Start() {
	var err error
	m.waitGroup.Add(1)

	go func(shutdownChannel chan struct{}, wg *sync.WaitGroup) {
		m.log.Info().Msg("Manager starting ...")
		m.client = pubsub.NewClient(m.log, "server")
		// connect to server websocket
		origin := "http://localhost/"
		url := "ws://localhost:12345/connect"
		err = m.client.Dial(url, origin)
		if err != nil {
			m.err <- err
		}
		err = m.client.Register("com.jtbonhomme.pubsub.general")
		if err != nil {
			m.err <- err
		}
		err = m.client.Register("com.jtbonhomme.pubsub.game")
		if err != nil {
			m.err <- err
		}

		m.started = true

		for {
			select {
			case <-shutdownChannel:
				m.log.Debug().Msg("shutdown manager websocket client")
				m.client.Shutdown()
				m.log.Debug().Msg("shutdown manager goroutine")
				wg.Done()
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
	m.log.Debug().Msg("close shutdown channel ...")
	close(m.shutdownChannel)
	m.log.Debug().Msg("wait from waitgroups being all done ...")
	m.waitGroup.Wait() // wait for all goroutines
	m.log.Debug().Msg("manager is shut down")
}
