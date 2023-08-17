package games

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/state"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
	"github.com/rs/zerolog"
)

type Game struct {
	log        *zerolog.Logger
	id         uuid.UUID `json:"id"`
	started    bool      `json:"started"`
	players    []string  `json:"players"`
	minplayers int       `json:"minplayers"`
	maxplayers int       `json:"maxplayers"`
	startTime  time.Time `json:"startTime"`
	endTime    time.Time `json:"endTime"`
	state      *state.State
}

// New creates a new game object with a minimum number of players
// required to join the game to be able to start it, and a maximum
// number of players who can join the game.
func New(l *zerolog.Logger, min, max int) *Game {
	gameID := uuid.New()

	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[game] %s", i) },
	}
	log := l.Output(output)

	g := Game{
		log:        &log,
		id:         gameID,
		minplayers: min,
		maxplayers: max,
		started:    false,
		players:    []string{},
		state:      state.New(l),
	}

	return &g
}

// Start starts the game. If the game is already started, or
// if the minimum player number registered is not reached, an
// error is returned.
func (game *Game) Start() error {
	if game.started {
		return fmt.Errorf("game already started")
	}

	players := game.players
	if game.minplayers != 0 && len(players) < game.minplayers {
		return fmt.Errorf("min player number %d not reached yet", game.minplayers)
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

// ID returns game's ID.
func (game *Game) ID() uuid.UUID {
	return game.id
}

// AddPlayer register a player to the game. If the player is already registered,
// the method does nothing. If the maximum number of players is alreary
// reached, of if the game is already started, the methods returns an error.
func (game *Game) AddPlayer(id string) error {
	if game.started {
		return errors.New("game alreay started")
	}

	if len(game.players) == game.maxplayers {
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
