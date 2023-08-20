package games

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	centrifuge "github.com/centrifugal/centrifuge-go"
	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

const (
	GameTopicPrefix                string = "game-"
	DefaultClientConnectionTimeout        = 10 * time.Second
)

type Game struct {
	log        *zerolog.Logger
	ID         uuid.UUID `json:"id"`
	players    []string
	MinPlayers int `json:"minPlayers"`
	MaxPlayers int `json:"maxPlayers"`
	startTime  time.Time
	endTime    time.Time
	started    bool
	TopicName  string `json:"topicName"`
	Name       string
	client     *centrifuge.Client
	sub        *centrifuge.Subscription
}

// New creates a new game object with a minimum number of players
// required to join the game to be able to start it, and a maximum
// number of players who can join the game.
func New(l *zerolog.Logger, min, max int) *Game {
	gameID := uuid.New()
	nameGenerator := namegenerator.NewNameGenerator(time.Now().UTC().UnixNano())
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[game] %s", i) },
	}
	log := l.Output(output)

	name := nameGenerator.Generate()

	g := Game{
		log:        &log,
		ID:         gameID,
		MinPlayers: min,
		MaxPlayers: max,
		players:    []string{},
		started:    false,
		TopicName:  GameTopicPrefix + name,
		Name:       name,
	}

	return &g
}

func (game *Game) messageHandler(e centrifuge.MessageEvent) {
	game.log.Info().Msgf("[%s] message received from server %s", game.Name, (e.Data))
}

func (game *Game) serverPublicationHandler(e centrifuge.ServerPublicationEvent) {
	game.log.Info().Msgf("[%s] server publication received %s", game.Name, (e.Data))
}

func (game *Game) publicationHandler(e centrifuge.PublicationEvent) {
	game.log.Info().Msgf("[%s] publication received %s", game.Name, (e.Data))
}

// Connect connects the game to the wrbsocket server.
func (game *Game) Connect() error {
	var err error
	var wg sync.WaitGroup
	var done = make(chan struct{})

	game.client = utils.NewClient(game.log,
		utils.DefaultWebsocketURL,
		&wg,
		utils.WithMessageHandler(game.messageHandler),
		utils.WithServerPublicationHandler(game.serverPublicationHandler),
	)

	wg.Add(1)
	err = game.client.Connect()
	if err != nil {
		return fmt.Errorf("[%s] centrifuge client connection error: %s", game.Name, err.Error())
	}

	// Will automatically close done channel and end select after all RPC have been called.
	go func() {
		game.log.Info().Msg("waiting client to connect...")
		defer close(done)
		wg.Wait()
	}()

	// Blocks until done channel is closed, or timeout occurs.
	select {
	case <-done:
		game.sub, err = utils.Subscribe(game.log,
			game.client,
			game.TopicName,
			utils.WithSubscriptionConfig(
				centrifuge.SubscriptionConfig{
					Data: []byte(`{"id":"` + game.ID.String() + `"}`),
				},
			),
			utils.WithPublicationHandler(game.publicationHandler),
		)
		if err != nil {
			return fmt.Errorf("[%s] subscription failed: %s", game.Name, err.Error())
		}

		return nil
	case <-time.After(DefaultClientConnectionTimeout):
		return fmt.Errorf("[%s] waiting for client connection failed with timeout", game.Name)
	}

}

// Start starts the game. If the game is already started, or
// if the minimum player number registered is not reached, an
// error is returned.
func (game *Game) Start() error {
	var err error

	if game.started {
		return fmt.Errorf("[%s] game already started", game.Name)
	}

	players := game.players
	if game.MinPlayers != 0 && len(players) < game.MinPlayers {
		return fmt.Errorf("[%s] min player number %d not reached yet", game.Name, game.MinPlayers)
	}

	game.started = true
	game.startTime = time.Now()

	// publication to all clients who subscribed to a channel
	_, err = game.publish(`{"type": "rpc", "emitter": "game", "id": "` + game.ID.String() + `", "data": "revealTwoCards"}`)
	if err != nil {
		game.log.Error().Msgf("[%s] publication error: %s", game.Name, err.Error())
	}
	return nil
}

// Publish sends a message on the game dedicated topic.
// An error is sent in case game is not connected to the internal centrifuge node.
func (game *Game) publish(message string) (centrifuge.PublishResult, error) {
	if game.sub == nil {
		return centrifuge.PublishResult{}, fmt.Errorf("subscription is nil")
	}

	res, err := game.sub.Publish(context.Background(),
		[]byte(message))
	if err != nil {
		return centrifuge.PublishResult{}, err
	}

	return res, nil
}

// Stop stops a started game. If the game is not started, an
// error is returned.
func (game *Game) Stop() error {
	if !game.started {
		return fmt.Errorf("[%s] game not started", game.Name)
	}

	game.started = false
	game.endTime = time.Now()

	return nil
}

// IsStarted returns true if the game is started.
func (game *Game) IsStarted() bool {
	return game.started
}

// AddPlayer register a player to the game. If the player is already registered,
// the method does nothing. If the maximum number of players is alreary
// reached, of if the game is already started, the methods returns an error.
func (game *Game) AddPlayer(id string) error {
	if game.started {
		return fmt.Errorf("[%s] game alreay started", game.Name)
	}

	if len(game.players) == game.MaxPlayers {
		return fmt.Errorf("[%s] maximum number of players alreay reached", game.Name)
	}

	if utils.ContainsString(game.players, id) {
		return fmt.Errorf("[%s] player id %s alreay joined the game", game.Name, id)
	}

	game.players = append(game.players, id)
	return nil
}

// Players returns game's registered players.
func (game *Game) Players() []string {
	return game.players
}

// PlayerInit -.
func (game *Game) PlayerInit() error {
	if !game.started {
		return fmt.Errorf("[%s] game not started", game.Name)
	}
	return nil
}
