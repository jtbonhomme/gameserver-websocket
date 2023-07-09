package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/pubsub"

	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
)

const skipFrameCount = 3

func main() {
	// Init logger
	var l zerolog.Level

	switch strings.ToLower("debug") {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(output).With().Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger()

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
