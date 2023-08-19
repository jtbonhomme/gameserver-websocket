package games

import (
	"errors"
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
		return fmt.Errorf("centrifuge client connection error: %s", err.Error())
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
		err = utils.Subscribe(game.log,
			game.client,
			game.TopicName,
			utils.WithSubscriptionConfig(
				centrifuge.SubscriptionConfig{
					Data: []byte(`{"id":"` + game.ID.String() + `"}`),
				},
			),
			utils.WithPublicationHandler(game.publicationHandler),
		)

		return nil
	case <-time.After(DefaultClientConnectionTimeout):
		return fmt.Errorf("waiting for client connection failed with timeout")
	}

}

// Start starts the game. If the game is already started, or
// if the minimum player number registered is not reached, an
// error is returned.
func (game *Game) Start() error {
	if game.started {
		return fmt.Errorf("game already started")
	}

	players := game.players
	if game.MinPlayers != 0 && len(players) < game.MinPlayers {
		return fmt.Errorf("min player number %d not reached yet", game.MinPlayers)
	}

	game.started = true
	game.startTime = time.Now()

	return nil
}

// Stop stops a started game. If the game is not started, an
// error is returned.
func (game *Game) Stop() error {
	if !game.started {
		return fmt.Errorf("game not started")
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
		return errors.New("game alreay started")
	}

	if len(game.players) == game.MaxPlayers {
		return errors.New("maximum number of players alreay reached")
	}

	if utils.ContainsString(game.players, id) {
		return fmt.Errorf("player id %s alreay joined the game", id)
	}

	game.players = append(game.players, id)
	return nil
}

// Players returns game's registered players.
func (game *Game) Players() []string {
	return game.players
}

// PlayerInit -.
func (game *Game) PlayerInit() {

}
