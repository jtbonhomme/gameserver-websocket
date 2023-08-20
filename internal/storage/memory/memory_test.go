package memory_test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
	"github.com/jtbonhomme/gameserver-websocket/internal/storage/memory"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
	"github.com/rs/zerolog"
)

var mgr *manager.Manager

type Response struct {
	Status string `json:"status"`
	Result string `json:"result"`
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
	c := utils.NewClient(&log, utils.DefaultWebsocketURL, &wg)
	err = c.Connect()
	if err != nil {
		log.Panic().Msgf("connect error: %s", err.Error())
	}

	log.Info().Msg("waiting client to connect...")
	wg.Wait()
	log.Info().Msg("client connected")

	// run tests suite
	exitVal := m.Run()
	c.Close()

	os.Exit(exitVal)
}
