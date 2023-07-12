package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/pubsub"
)

type RideMessage struct {
	RideID         string  `json:"ride_id"`
	PointIdx       int     `json:"point_idx"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	Timestamp      string  `json:"timestamp"`
	MeterReading   float64 `json:"meter_reading"`
	MeterIncrement float64 `json:"meter_increment"`
	RideStatus     string  `json:"ride_status"`
	PassengerCount int     `json:"passenger_count"`
}

const tsFormat string = "2006-01-02T15:04:05.0000-07:00"

var rideStatus = []string{"STARTED", "WAITING", "FINISHED", "BLOCKED"}

func receiveMessageHandler(payload []byte) error {
	m := pubsub.Message{}

	err := json.Unmarshal(payload, &m)
	if err != nil {
		return fmt.Errorf("error unmarshaling payload: %w", err)
	}
	fmt.Printf("received message/ %#v\n", m)
	return nil
}

func main() {
	var err error

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[client] %s", i) },
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	c := pubsub.NewClient(&logger, "main-pubsub-client")
	// connect to server websocket
	origin := "http://localhost/"
	url := "ws://localhost:12345/connect"
	err = c.Dial(url, origin)
	if err != nil {
		logger.Fatal().Err(err).Msg("error dialing websocket server")
	}
	err = c.Register("com.jtbonhomme.pubsub.general")
	if err != nil {
		logger.Fatal().Err(err).Msg("error registering to topic")
	}

	// send a message
	payload, err := json.Marshal(RideMessage{
		RideID:         uuid.NewString(),
		PointIdx:       r.Intn(10),
		Latitude:       r.Float64() + float64(48),
		Longitude:      r.Float64() + float64(2),
		Timestamp:      time.Now().Format(tsFormat),
		MeterReading:   r.Float64() * 50,
		MeterIncrement: r.Float64() * 2,
		PassengerCount: r.Intn(5),
		RideStatus:     rideStatus[r.Intn(3)],
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Publish("com.jtbonhomme.pubsub.game", payload)

	c.Read(receiveMessageHandler)
}
