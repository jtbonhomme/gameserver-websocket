package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
)

func main() {
	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[main] %s", i) },
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	logger.Info().Msg("start manager")
	mgr := manager.New(&logger)
	err := mgr.Start()
	if err != nil {
		logger.Panic().Msgf("manager start error: %w", err)
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
