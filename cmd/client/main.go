package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/goombaio/namegenerator"
	"github.com/jtbonhomme/gameserver-websocket/internal/games"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/rs/zerolog"
)

var wg sync.WaitGroup

const (
	serverTopicName string = "server-general"
)

type Response struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

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
		log.Info().Msg("Connected")
		wg.Done()
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

	wg.Add(1)
	err = c.Connect()
	if err != nil {
		log.Panic().Msgf("connect error: %s", err.Error())
	}

	log.Info().Msg("waiting client to connect...")
	wg.Wait()
	log.Info().Msg("client connected")

	serverTopic, err := c.NewSubscription("server-general")
	if err != nil {
		log.Error().Msgf("subscription to %s creation error: %s", serverTopicName, err.Error())
	}

	serverTopic.OnJoin(func(e centrifuge.JoinEvent) {
		log.Info().Msgf("[%s] join event: %s", serverTopicName, e.ClientInfo.Client)
	})

	serverTopic.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Info().Msgf("[%s] subscription error event: %s", serverTopicName, e.Error.Error())
	})

	serverTopic.OnPublication(func(e centrifuge.PublicationEvent) {
		log.Info().Msgf("[%s] publication event: %s", serverTopicName, string(e.Data))
	})

	serverTopic.OnSubscribing(func(e centrifuge.SubscribingEvent) {
		log.Info().Msgf("[%s] subscribing event: %s", serverTopicName, e.Reason)
	})

	serverTopic.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		log.Info().Msgf("[%s] subscribed event", serverTopicName)
	})

	err = serverTopic.Subscribe()
	if err != nil {
		log.Error().Msgf("subscription error: %s", err.Error())
	}

	var result centrifuge.RPCResult
	nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())

	clientName := nameGenerator.Generate()
	log.Info().Msgf("generated client name: %s", clientName)

	result, err = c.RPC(context.Background(), "registerPlayer", []byte(`{"name":"`+clientName+`"}`))
	if err != nil {
		log.Panic().Msgf("error executing RPC: %s", err.Error())
	}
	log.Debug().Msgf("registerPlayer result: %s", string(result.Data))

	var response Response
	err = json.Unmarshal(result.Data, &response)
	if err != nil {
		log.Panic().Msgf("error unmarshaling Player: %s", err.Error())
	}
	log.Debug().Msgf("response %#v", response)

	var player players.Player
	err = json.Unmarshal([]byte(response.Result), &player)
	if err != nil {
		log.Panic().Msgf("error unmarshaling Player: %s", err.Error())
	}
	log.Debug().Msgf("player %#v", player)

	result, err = c.RPC(context.Background(), "createGame", []byte(`{"minPlayers": 1, "maxPlayers": 2}`))
	if err != nil {
		log.Panic().Msgf("error executing RPC: %s", err.Error())
	}
	log.Debug().Msgf("createGame result: %s", string(result.Data))

	err = json.Unmarshal(result.Data, &response)
	if err != nil {
		log.Panic().Msgf("error unmarshaling Player: %s", err.Error())
	}
	log.Debug().Msgf("response %#v", response)

	var game games.Game
	err = json.Unmarshal([]byte(response.Result), &game)
	if err != nil {
		log.Panic().Msgf("error unmarshaling Player: %s", err.Error())
	}
	log.Debug().Msgf("game %#v", game)

	result, err = c.RPC(context.Background(), "joinGame", []byte(`{"idGame": "`+game.ID.String()+`", "idPlayer": "`+player.ID.String()+`"}`))
	if err != nil {
		log.Panic().Msgf("error executing RPC: %s", err.Error())
	}
	log.Debug().Msgf("joinGame result: %s", string(result.Data))

	gameTopic, err := c.NewSubscription(game.TopicName)
	if err != nil {
		log.Error().Msgf("subscription to %s creation error: %s", game.TopicName, err.Error())
	}

	gameTopic.OnJoin(func(e centrifuge.JoinEvent) {
		log.Info().Msgf("[%s] join event: %s", game.TopicName, e.ClientInfo.Client)
	})

	gameTopic.OnError(func(e centrifuge.SubscriptionErrorEvent) {
		log.Info().Msgf("[%s] subscription error event: %s", game.TopicName, e.Error.Error())
	})

	gameTopic.OnPublication(func(e centrifuge.PublicationEvent) {
		log.Info().Msgf("[%s] publication event: %s", game.TopicName, string(e.Data))
	})

	gameTopic.OnSubscribing(func(e centrifuge.SubscribingEvent) {
		log.Info().Msgf("[%s] subscribing event: %s", game.TopicName, e.Reason)
	})

	gameTopic.OnSubscribed(func(e centrifuge.SubscribedEvent) {
		log.Info().Msgf("[%s] subscribed event", game.TopicName)
	})

	err = gameTopic.Subscribe()
	if err != nil {
		log.Error().Msgf("subscription error: %s", err.Error())
	}

	result, err = c.RPC(context.Background(), "startGame", []byte(`{"id": "`+game.ID.String()+`"}`))
	if err != nil {
		log.Panic().Msgf("error executing RPC: %s", err.Error())
	}
	log.Debug().Msgf("startGame result: %s", string(result.Data))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Info().Msg("received signal: " + s.String())
	// Shutdown

	log.Info().Msg("exit")
}
