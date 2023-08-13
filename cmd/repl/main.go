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
	wg.Add(1)
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

	err = c.Connect()
	if err != nil {
		log.Error().Msgf("connect error: %s", err.Error())
		return
	}
	wg.Wait()

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
