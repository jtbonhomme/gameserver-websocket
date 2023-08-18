package games

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/goombaio/namegenerator"
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/state"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

const (
	GameTopicPrefix string = "game-"
)

type Game struct {
	log        *zerolog.Logger
	ID         uuid.UUID `json:"id"`
	started    bool
	players    []string
	MinPlayers int `json:"minPlayers"`
	MaxPlayers int `json:"maxPlayers"`
	startTime  time.Time
	endTime    time.Time
	state      *state.State
	TopicName  string `json:"topicName"`
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

	g := Game{
		log:        &log,
		ID:         gameID,
		MinPlayers: min,
		MaxPlayers: max,
		started:    false,
		players:    []string{},
		state:      state.New(l),
		TopicName:  GameTopicPrefix + nameGenerator.Generate(),
	}

	return &g
}

// Start starts the game. If the game is already started, or
// if the minimum player number registered is not reached, an
// error is returned.
func (game *Game) Start() error {
	var err error

	if game.started {
		return fmt.Errorf("game already started")
	}

	players := game.players
	if game.MinPlayers != 0 && len(players) < game.MinPlayers {
		return fmt.Errorf("min player number %d not reached yet", game.MinPlayers)
	}

	err = game.state.Start()
	if err != nil {
		return fmt.Errorf("game state can not be updated to start: %s", err.Error())
	}

	game.started = true
	game.startTime = time.Now()

	return nil
}

// Stop stops a started game. If the game is not started, an
// error is returned.
func (game *Game) Stop() error {
	var err error

	if !game.started {
		return fmt.Errorf("game not started")
	}

	err = game.state.Stop()
	if err != nil {
		return fmt.Errorf("game state can not be updated to stop: %s", err.Error())
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

// CurrentState returns internal game state.
func (game *Game) CurrentState() string {
	return game.state.Current()
}
