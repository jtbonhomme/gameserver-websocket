package manager

import (
	"fmt"
	"os"
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
	shutdownChannel chan bool // Channel used to stop manager.
	waitGroup       *sync.WaitGroup
	c               *pubsub.Client
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
		shutdownChannel: make(chan bool),
		err:             make(chan error),
		waitGroup:       &sync.WaitGroup{},
	}
}

func (m *Manager) listen(msg []byte) error {
	m.log.Info().Msgf("received message %s", string(msg))
	return nil
}

func (m *Manager) Start() {
	m.c = pubsub.NewClient(m.log, "manager-pubsub-client")
	var err error
	m.log.Info().Msg("starting ...")

	// connect to server websocket
	origin := "http://localhost/"
	url := "ws://localhost:12345/connect"
	err = m.c.Dial(url, origin)
	if err != nil {
		m.err <- err
	}
	err = m.c.Register("com.jtbonhomme.pubsub.general")
	if err != nil {
		m.err <- err
	}
	err = m.c.Register("com.jtbonhomme.pubsub.game")
	if err != nil {
		m.err <- err
	}

	m.started = true

	m.c.Read(m.listen)

	m.waitGroup.Add(1)
	go func(shutdownChannel chan bool, wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			select {
			case <-shutdownChannel:
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
	m.started = false
	m.log.Info().Msg("shuting down ...")
	m.shutdownChannel <- true
	m.log.Info().Msgf("waiting for go routine to stop ...")
	m.waitGroup.Wait()
	m.log.Info().Msgf("all go routine stopped")
	close(m.shutdownChannel)
	m.c.Shutdown()
	m.log.Info().Msgf("stopped")
}
