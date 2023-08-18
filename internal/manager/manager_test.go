package manager_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
	"github.com/jtbonhomme/gameserver-websocket/internal/storage/memory"
	"github.com/rs/zerolog"
)

var mgr *manager.Manager

type Response struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func newClient(log *zerolog.Logger, wg *sync.WaitGroup) *centrifuge.Client {
	wsURL := "ws://localhost:8000/connection/websocket"

	c := centrifuge.NewJsonClient(wsURL, centrifuge.Config{
		Name:    "listening-go-client",
		Version: "0.0.1",
	})

	c.OnConnecting(func(_ centrifuge.ConnectingEvent) {
		log.Info().Msg("Connecting")
	})

	c.OnConnected(func(_ centrifuge.ConnectedEvent) {
		log.Info().Msg("Connected")
		wg.Done()
	})

	c.OnDisconnected(func(e centrifuge.DisconnectedEvent) {
		log.Info().Msgf("Disconnected event: %d %s", e.Code, e.Reason)
		// TODO automatic reconnect?
	})

	c.OnError(func(e centrifuge.ErrorEvent) {
		log.Info().Msgf("error: %s", e.Error.Error())
	})

	c.OnMessage(func(e centrifuge.MessageEvent) {
		log.Info().Msgf("Message received from server %s", string(e.Data))
	})

	return c
}

func newLogger() zerolog.Logger {
	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[client] %s", i) },
	}
	log := zerolog.New(output).With().Timestamp().Logger()

	return log
}

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	var err error

	log := newLogger()
	log.Info().Msgf("Do stuff BEFORE the tests!")
	log.Info().Msg("start client")

	// concrete memory test storage implementation
	s := memory.New(&log)
	mgr = manager.New(&log, s)
	err = mgr.Start()
	if err != nil {
		log.Err(err).Msg("error starting manager")
		os.Exit(1)
	}

	// start client to receive publications
	wg.Add(1)
	c := newClient(&log, &wg)
	err = c.Connect()
	if err != nil {
		log.Panic().Msgf("connect error: %s", err.Error())
	}

	log.Info().Msg("waiting client to connect...")
	wg.Wait()
	log.Info().Msg("client connected")

	// run tests suite
	exitVal := m.Run()
	log.Info().Msgf("Do stuff AFTER the tests!")
	c.Close()

	os.Exit(exitVal)
}
