package games_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jtbonhomme/gameserver-websocket/internal/games"
	"github.com/rs/zerolog"
)

func TestGame_Start(t *testing.T) {
	var err error

	logger := zerolog.Nop()

	game := games.New(&logger, 2, 3)

	err = game.Start()
	if err == nil {
		t.Errorf("expected an error when starting the game as required numberof player is not reached: %v", err)
	}

	playerID := "player1"
	err = game.AddPlayer(playerID)
	if err != nil {
		t.Errorf("Unexpected error when adding a player: %v", err)
	}

	if len(game.Players()) != 1 {
		t.Errorf("Expected 1 player, but got %d", len(game.Players()))
	}

	// Add the same player again, should not add a duplicate
	err = game.AddPlayer(playerID)
	if err == nil {
		t.Errorf("expected error when adding same player twice: %v", err)
	}

	if len(game.Players()) != 1 {
		t.Errorf("Expected 1 player after adding duplicate, but got %d", len(game.Players()))
	}

	// Add another player
	err = game.AddPlayer("player2")
	if err != nil {
		t.Errorf("Unexpected error when adding a player: %v", err)
	}
	if len(game.Players()) != 2 {
		t.Errorf("Expected 2 players, but got %d", len(game.Players()))
	}

	// Start the game successfully
	err = game.Start()
	if err != nil {
		t.Errorf("Unexpected error when starting the game: %v", err)
	}

	if !game.IsStarted() {
		t.Error("Expected game to be started, but it's not")
	}

	// Attempt to start an already started game, should return an error
	err = game.Start()
	if err == nil {
		t.Error("Expected error when starting an already started game")
	}
}

func TestGame_Stop(t *testing.T) {
	var err error

	logger := zerolog.Nop()

	game := games.New(&logger, 1, 2)

	err = game.Start()
	if err == nil {
		t.Errorf("expected an error when starting the game as required numberof player is not reached: %v", err)
	}

	playerID := "player1"
	err = game.AddPlayer(playerID)
	if err != nil {
		t.Errorf("Unexpected error when adding a player: %v", err)
	}

	if len(game.Players()) != 1 {
		t.Errorf("Expected 1 player, but got %d", len(game.Players()))
	}

	// Add the same player again, should not add a duplicate
	err = game.AddPlayer(playerID)
	if err == nil {
		t.Errorf("expected error when adding a player already added: %v", err)
	}
	if len(game.Players()) != 1 {
		t.Errorf("Expected 1 player after adding duplicate, but got %d", len(game.Players()))
	}

	// Add another player
	err = game.AddPlayer("player2")
	if err != nil {
		t.Errorf("Unexpected error when adding a player: %v", err)
	}
	if len(game.Players()) != 2 {
		t.Errorf("Expected 2 players, but got %d", len(game.Players()))
	}

	// Add another player
	err = game.AddPlayer("player3")
	if err == nil {
		t.Errorf("expected error when adding a player while max required number of players is reached: %v", err)
	}

	if len(game.Players()) != 2 {
		t.Errorf("Expected 2 players, but got %d", len(game.Players()))
	}

	// Start the game successfully
	err = game.Start()
	if err != nil {
		t.Errorf("Unexpected error when starting the game: %v", err)
	}

	if !game.IsStarted() {
		t.Error("Expected game to be started, but it's not")
	}

	// Stop the game successfully
	err = game.Stop()
	if err != nil {
		t.Errorf("Unexpected error when stopping the game: %v", err)
	}

	if game.IsStarted() {
		t.Error("Expected game to be stopped, but it's still started")
	}

	// Attempt to stop an already stopped game, should return an error
	err = game.Stop()
	if err == nil {
		t.Error("Expected error when stopping an already stopped game")
	}
}

func TestGame_ID(t *testing.T) {
	logger := zerolog.Nop()

	game := games.New(&logger, 2, 4)

	id := game.ID()
	if id == uuid.Nil {
		t.Error("Expected non-nil game ID")
	}
}
