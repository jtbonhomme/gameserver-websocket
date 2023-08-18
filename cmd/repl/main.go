package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

func newClient(log *zerolog.Logger, wg *sync.WaitGroup) *centrifuge.Client {
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
		wg.Done()
	})

	c.OnDisconnected(func(_ centrifuge.DisconnectedEvent) {
		log.Info().Msg("Disconnected")
	})

	c.OnError(func(e centrifuge.ErrorEvent) {
		log.Info().Msgf("on error: %s", e.Error.Error())
	})

	c.OnMessage(func(e centrifuge.MessageEvent) {
		log.Info().Msgf("Message received from server %s", string(e.Data))
	})

	return c
}

func main() {
	var err error
	var wg sync.WaitGroup

	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[client] %s", i) },
	}
	log := zerolog.New(output).With().Timestamp().Logger()
	log.Info().Msg("start client")
	c := newClient(&log, &wg)
	defer c.Close()

	wg.Add(1)
	err = c.Connect()
	if err != nil {
		log.Error().Msgf("connect error: %s", err.Error())
		return
	}

	log.Info().Msg("waiting client to connect...")
	wg.Wait()
	log.Info().Msg("client connected")

	serverTopic, err := c.NewSubscription("server-general")
	if err != nil {
		log.Error().Msgf("subscription creation error: %s", err.Error())
	}

	serverTopic.OnJoin(func(e centrifuge.JoinEvent) {
		log.Info().Msgf("[server-general] join event: %s", e.ClientInfo.Client)
	})

	serverTopic.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Info().Msgf("[server-general] subscription error event: %s", e.Error.Error())
	})

	serverTopic.OnPublication(func(e centrifuge.PublicationEvent) {
		log.Info().Msgf("[server-general] publication event: %s", string(e.Data))
	})

	serverTopic.OnSubscribing(func(e centrifuge.SubscribingEvent) {
		log.Info().Msgf("[server-general] subscribing event: %s", e.Reason)
	})

	serverTopic.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		log.Info().Msgf("[server-general] subscribed event")
		wg.Done()
	})

	wg.Add(1)
	err = serverTopic.Subscribe()
	if err != nil {
		log.Error().Msgf("subscription error: %s", err.Error())
	}

	log.Info().Msg("waiting client to subscribe...")
	wg.Wait()
	log.Info().Msg("client subscribed")

	rpc := replRun()
	validate := validator.New()
	err = validate.Var(rpc.Payload, "json")
	if err != nil {
		log.Error().Msgf("error: string \"%s\" is not a valid JSON", rpc.Payload)
		return
	}

	result, err := c.RPC(context.Background(), rpc.Method, []byte(rpc.Payload))
	if err != nil {
		log.Error().Msgf("error executing RPC: %s", err.Error())
	}
	log.Printf("RPC result: %s", string(result.Data))

	log.Info().Msg("exit")
}
