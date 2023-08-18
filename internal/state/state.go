package state

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/looplab/fsm"
)

const (
	// States
	CreatedState string = "created"
	StartedState string = "started"
	StoppedState string = "stopped"
	// Events
	StartEvent string = "start"
	StopEvent  string = "stop"
)

// State structure maintains game state.
type State struct {
	state *fsm.FSM
	log   *zerolog.Logger
}

// New initializes a new game State.
func New(l *zerolog.Logger) *State {
	events := fsm.Events{
		{
			Name: StartEvent,
			Src:  []string{CreatedState},
			Dst:  StartedState,
		},
		{
			Name: StopEvent,
			Src:  []string{StartedState},
			Dst:  StoppedState,
		},
	}

	cb := fsm.Callbacks{}

	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[state] %s", i) },
	}
	log := l.Output(output)

	return &State{
		state: fsm.NewFSM(
			CreatedState,
			events,
			cb,
		),
		log: &log,
	}
}

// Current returns internal game state.
func (s *State) Current() string {
	return s.state.Current()
}

// Start sends "start" event to game engine internal state.
func (s *State) Start() error {
	return s.state.Event(context.Background(), StartEvent)
}

// Stop sends "stop" event to game engine internal state.
func (s *State) Stop() error {
	return s.state.Event(context.Background(), StopEvent)
}
