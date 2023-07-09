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
	"github.com/jtbonhomme/gameserver-websocket/internal/websocket"
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
	broker := pubsub.New(&logger)
	srv := websocket.New(&logger, broker)
	srv.Start()

	mgr := manager.New(&logger)
	mgr.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info().Msg("signal: " + s.String())
		// Shutdown
		srv.Shutdown()
		mgr.Shutdown()
	}

	logger.Info().Msg("Exit")
}
