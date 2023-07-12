package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/pubsub"

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

	// 1. Server Setup
	logger.Info().Msg("start broker")
	broker := pubsub.New(&logger)
	broker.Start()

	// todo: use websocket to send logs: provide an io.Writer implementing object
	// todo: manager websocket availability (because it starts in a goroutine)

	logger.Info().Msg("start manager")
	mgr := manager.New(&logger)
	mgr.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info().Msg("received signal: " + s.String())
		// Shutdown
		logger.Info().Msg("start shutdown procedure")
		broker.Shutdown()
		mgr.Shutdown()
	case err := <-broker.Error():
		logger.Err(err).Msg("broker error")
	case err := <-mgr.Error():
		logger.Err(err).Msg("manager error")
	}
	logger.Info().Msg("exit")
}
