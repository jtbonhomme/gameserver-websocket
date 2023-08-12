package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
	"github.com/jtbonhomme/gameserver-websocket/internal/storage/memory"
)

func main() {
	var err error

	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[main] %s", i) },
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	m := memory.New()
	logger.Info().Msg("start manager")
	mgr := manager.New(&logger, m)
	err = mgr.Start()
	if err != nil {
		logger.Panic().Msgf("manager start error: %s", err.Error())
	}

	// todo: use websocket to send logs: provide an io.Writer implementing object
	// todo: manager websocket availability (because it starts in a goroutine)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info().Msg("received signal: " + s.String())
		// Shutdown
		logger.Info().Msg("start shutdown procedure")
		mgr.Shutdown()
	case err := <-mgr.Error():
		logger.Err(err).Msg("manager error")
	}
	logger.Info().Msg("exit")
}
