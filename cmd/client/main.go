package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/rs/zerolog"
)

func newClient(log *zerolog.Logger) *centrifuge.Client {
	wsURL := "ws://localhost:8000/connection/websocket"
	c := centrifuge.NewJsonClient(wsURL, centrifuge.Config{
		Name:    "centrifuge-go",
		Data:    []byte(`{"name":"totoro"}`),
		Version: "0.0.1",
	})

	c.OnConnecting(func(_ centrifuge.ConnectingEvent) {
		log.Info().Msg("Connecting")
	})
	c.OnConnected(func(_ centrifuge.ConnectedEvent) {
		log.Info().Msg("Connected")
	})
	c.OnDisconnected(func(_ centrifuge.DisconnectedEvent) {
		log.Info().Msg("Disconnected")
	})
	c.OnError(func(e centrifuge.ErrorEvent) {
		log.Info().Msgf("error: %s", e.Error.Error())
	})
	c.OnMessage(func(e centrifuge.MessageEvent) {
		log.Info().Msgf("Message received from server %s", string(e.Data))

		// When issue blocking requests from inside event handler we must use
		// a goroutine. Otherwise, connection read loop will be blocked.
		go func() {
			result, err := c.RPC(context.Background(), "method", []byte(`{"action":"eat"}`))
			if err != nil {
				log.Info().Msgf("%s", err.Error())
				return
			}
			log.Printf("RPC result 2: %s", string(result.Data))
		}()
	})
	return c
}

func main() {
	var err error

	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[client] %s", i) },
	}
	log := zerolog.New(output).With().Timestamp().Logger()
	log.Info().Msg("start client")
	c := newClient(&log)
	defer c.Close()

	err = c.Connect()
	if err != nil {
		log.Panic().Msgf("connect error: %s", err.Error())
	}

	result, err := c.RPC(context.Background(), "method", []byte(`{"action":"drink"}`))
	if err != nil {
		log.Panic().Msgf("rpc error: %s", err.Error())
	}

	log.Printf("RPC result: %s", string(result.Data))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Msg("received signal: " + s.String())
		// Shutdown
		log.Info().Msg("start shutdown procedure")
	}

	log.Info().Msg("exit")
}
