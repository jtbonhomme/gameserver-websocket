package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/pubsub"

	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
)

const skipFrameCount = 3

func main() {
	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	// 1. Server Setup
	logger.Info().Msg("Server: start broker")
	broker := pubsub.New(&logger)
	broker.Start()

	// todo: manager websocket availability (because it starts in a goroutine)

	logger.Info().Msg("Server: start manager")
	mgr := manager.New(&logger)
	mgr.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info().Msg("signal: " + s.String())
		// Shutdown
		broker.Shutdown()
		mgr.Shutdown()
	}
	logger.Info().Msg("Exit")
}
