package memory_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/storage/memory"
)

func TestMemoryPlayer(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[main] %s", i) },
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	// concrete memory test storage implementation
	s := memory.New(&logger)

	// first time player registration
	player1, err := s.Register("", "name1")
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}
	if player1.ID.String() == "" {
		t.Errorf("error, expected new UUID created for player name1")
	}
	id := player1.ID.String()

	player2, err := s.Register("", "name2")
	if err != nil {
		t.Errorf("error while registering name2: %s", err.Error())
	}
	if player2.ID.String() == "" {
		t.Errorf("error, expected new UUID created for player name2")
	}

	// register already recorded player (name1)
	player3, err := s.Register(id, "name1")
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}
	if player3.ID.String() == "" {
		t.Errorf("error, expected new UUID created for player name1")
	}
	if player3.ID.String() != id {
		t.Errorf("error, expected new UUID created for player name1")
	}

	players := s.ListAll()
	if len(players) != 2 {
		t.Errorf("expected 2 registered players and got %d", len(players))
	}

	// unregister unknown player
	err = s.Unregister("fake-id")
	if err == nil {
		t.Errorf("expected error while registering unknown player but got nil")
	}

	players = s.ListAll()
	if len(players) != 2 {
		t.Errorf("expected 2 registered players and got %d", len(players))
	}

	// unregister already recorded player (name1)
	err = s.Unregister(id)
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}

	players = s.ListAll()
	if len(players) != 1 {
		t.Errorf("expected 1 registered players and got %d", len(players))
	}
	if players[0].ID != player2.ID {
		t.Errorf("expected only registered player to be %s (%s) and got %s (%s)", player2.Name, player2.ID, players[0].Name, players[0].ID)
	}
}
