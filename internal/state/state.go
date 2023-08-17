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
	CreatedState      string = "created"
	InitializingState string = "initializing"
	StartedState      string = "started"
	EndedState        string = "ended"
	// Events
	InitEvent  string = "init"
	StartEvent string = "start"
	EndEvent   string = "end"
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
			Name: InitEvent,
			Src:  []string{CreatedState},
			Dst:  InitializingState,
		},
		{
			Name: StartEvent,
			Src:  []string{InitializingState},
			Dst:  StartedState,
		},
		{
			Name: EndEvent,
			Src:  []string{StartedState},
			Dst:  EndedState,
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

// End sends "end" event to game engine internal state.
func (s *State) End() error {
	return s.state.Event(context.Background(), EndEvent)
}
