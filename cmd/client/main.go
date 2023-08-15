package main

import (
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
		Name:    "listening-go-client",
		Version: "0.0.1",
	})

	c.OnConnecting(func(_ centrifuge.ConnectingEvent) {
		log.Info().Msg("Connecting")
	})

	c.OnConnected(func(_ centrifuge.ConnectedEvent) {
		var err error
		log.Info().Msg("Connected")

		subServer, err := c.NewSubscription("com.jtbonhomme.server")
		if err != nil {
			log.Error().Msgf("subscription creation error: %s", err.Error())
		}

		subServer.OnJoin(func(e centrifuge.JoinEvent) {
			log.Info().Msgf("[com.jtbonhomme.server] join event: %s", e.ClientInfo.Client)
		})

		subServer.OnError(func(e centrifuge.SubscriptionErrorEvent) {
			log.Info().Msgf("[com.jtbonhomme.server] subscription error event: %s", e.Error.Error())
		})

		subServer.OnPublication(func(e centrifuge.PublicationEvent) {
			log.Info().Msgf("[com.jtbonhomme.server] publication event: %s", string(e.Data))
		})

		subServer.OnSubscribing(func(e centrifuge.SubscribingEvent) {
			log.Info().Msgf("[com.jtbonhomme.server] subscribing event: %s", e.Reason)
		})

		subServer.OnSubscribed(func(e centrifuge.SubscribedEvent) {
			log.Info().Msgf("[com.jtbonhomme.server] subscribed event")
		})

		err = subServer.Subscribe()
		if err != nil {
			log.Error().Msgf("subscription error: %s", err.Error())
		}
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

func main() {
	var err error

	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[client] %s", i) },
	}
	log := zerolog.New(output).With().Timestamp().Logger()
	log.Info().Msg("start client")
	c := newClient(&log)
	defer c.Close()

	err = c.Connect()
	if err != nil {
		log.Panic().Msgf("connect error: %s", err.Error())
	}

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
