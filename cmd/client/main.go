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
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/games"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

type Response struct {
	Status string `json:"status"`
	Result string `json:"result"`
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
	c := utils.NewClient(&log, utils.DefaultWebsocketURL, &wg)
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

	_, err = utils.Subscribe(&log, c, utils.ServerPublishChannel)
	if err != nil {
		log.Error().Msgf("subscribe error: %s", err.Error())
		return
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

	var game games.Game
	err = json.Unmarshal([]byte(response.Result), &game)
	if err != nil {
		log.Panic().Msgf("error unmarshaling Player: %s", err.Error())
	}

	log.Debug().Msgf("JOIN GAME %s", game.Name)
	result, err = c.RPC(context.Background(), "joinGame", []byte(`{"idGame": "`+game.ID.String()+`", "idPlayer": "`+player.ID.String()+`"}`))
	if err != nil {
		log.Panic().Msgf("error executing RPC: %s", err.Error())
	}
	log.Debug().Msgf("joinGame result: %s", string(result.Data))

	var publicationHandler = func(e centrifuge.PublicationEvent) {
		log.Info().Msgf("[%s] publication event: %s", game.TopicName, string(e.Data))
		var data struct {
			Type    string `json:"type"`
			Emitter string `json:"emitter"`
			ID      string `json:"id"`
			Data    string `json:"data"`
		}

		err = json.Unmarshal(e.Data, &data)
		if err != nil {
			log.Error().Msgf("error while unmarshaling data %q: %s", string(e.Data), err.Error())
		}
		log.Info().Msgf("[%s] publication event: %#v", game.TopicName, data)

		if data.Type == "rpc" && data.Data == "revealTwoCards" {
			log.Info().Msgf("[%s] revealTwoCards publication event !", game.TopicName)
			wg.Done()
		}
	}

	_, err = utils.Subscribe(&log,
		c,
		game.TopicName,
		utils.WithSubscriptionConfig(
			centrifuge.SubscriptionConfig{
				Data: []byte(`{"id":"` + player.ID.String() + `"}`),
			},
		),
		utils.WithPublicationHandler(publicationHandler),
	)
	if err != nil {
		log.Error().Msgf("subscribe error: %s", err.Error())
		return
	}

	log.Debug().Msgf("subscribed topic %s", game.TopicName)

	wg.Add(1)
	log.Debug().Msgf("START GAME %s", game.Name)
	result, err = c.RPC(context.Background(), "startGame", []byte(`{"id": "`+game.ID.String()+`"}`))
	if err != nil {
		log.Error().Msgf("error executing RPC: %s", err.Error())
		return
	}
	log.Debug().Msgf("startGame result: %s", string(result.Data))

	log.Debug().Msgf("waiting subscribe event revealTwoCards ...")
	wg.Wait()
	log.Debug().Msgf("received subscribe event revealTwoCards")

	log.Debug().Msgf("wait 2s to call rpc playerInit")
	time.Sleep(2 * time.Second)
	log.Debug().Msgf("PLAYERINIT GAME %s", game.Name)
	result, err = c.RPC(context.Background(), "playerInit", []byte(`{"idGame": "`+game.ID.String()+`", "idPlayer": "`+player.ID.String()+`"}`))
	if err != nil {
		log.Panic().Msgf("error executing RPC: %s", err.Error())
	}
	log.Debug().Msgf("playerInit result: %s", string(result.Data))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Info().Msg("received signal: " + s.String())
	// Shutdown

	log.Info().Msg("exit")
}
