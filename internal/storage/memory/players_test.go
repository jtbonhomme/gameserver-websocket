package memory_test

import (
	"testing"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/storage/memory"
)

func TestMemoryPlayer(t *testing.T) {
	logger := zerolog.Nop()

	// concrete memory test storage implementation
	mem := memory.New(&logger)

	// first time player registration
	player1, err := mem.RegisterPlayer("", "name1")
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}

	if player1 == nil {
		t.Fatal("Expected player1 to be created, got nil")
	}

	if player1 != nil && player1.ID.String() == "" {
		t.Errorf("error, expected new UUID created for player name1")
	}
	id := player1.ID.String()

	player2, err := mem.RegisterPlayer("", "name2")
	if err != nil {
		t.Errorf("error while registering name2: %s", err.Error())
	}
	if player2.ID.String() == "" {
		t.Errorf("error, expected new UUID created for player name2")
	}

	// register already recorded player (name1)
	player3, err := mem.RegisterPlayer(id, "name1")
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}
	if player3.ID.String() == "" {
		t.Errorf("error, expected new UUID created for player name1")
	}
	if player3.ID.String() != id {
		t.Errorf("error, expected new UUID created for player name1")
	}

	players := mem.ListPlayers()
	if len(players) != 2 {
		t.Errorf("expected 2 registered players and got %d", len(players))
	}

	// unregister unknown player
	err = mem.UnregisterPlayer("fake-id")
	if err == nil {
		t.Errorf("expected error while registering unknown player but got nil")
	}

	players = mem.ListPlayers()
	if len(players) != 2 {
		t.Errorf("expected 2 registered players and got %d", len(players))
	}

	// unregister already recorded player (name1)
	err = mem.UnregisterPlayer(id)
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}

	players = mem.ListPlayers()
	if len(players) != 1 {
		t.Errorf("expected 1 registered players and got %d", len(players))
	}
	if players[0].ID != player2.ID {
		t.Errorf("expected only registered player to be %s (%s) and got %s (%s)", player2.Name, player2.ID, players[0].Name, players[0].ID)
	}
}
