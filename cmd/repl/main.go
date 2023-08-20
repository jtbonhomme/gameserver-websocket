package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

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
