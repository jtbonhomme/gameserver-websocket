package memory_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/storage/memory"
)

func TestMemoryGamesErrors(t *testing.T) {
	var err error

	logger := zerolog.Nop()

	// concrete memory test storage implementation
	mem := memory.New(&logger)

	// Start unknown game with error
	err = mem.StartGame("fake")
	if err == nil {
		t.Errorf("Expected error when starting unknown game")
	}

	// Start nil game id with error
	err = mem.StartGame(uuid.Nil.String())
	if err == nil {
		t.Errorf("Expected error when starting unknown game")
	}

	// Stop unknown game with error
	err = mem.StopGame("fake")
	if err == nil {
		t.Errorf("Expected error when stopping unknown game")
	}

	// Stop nil game id with error
	err = mem.StopGame(uuid.Nil.String())
	if err == nil {
		t.Errorf("Expected error when stopping nil game id")
	}

	_, err = mem.IsGameStarted("fake")
	if err == nil {
		t.Errorf("Expected error when checking if unknown game is started")
	}

	_, err = mem.IsGameStarted(uuid.Nil.String())
	if err == nil {
		t.Errorf("Expected error when checking if nil game id is started")
	}

	err = mem.JoinGame("fake", uuid.NewString())
	if err == nil {
		t.Errorf("Expected error when joining unknown game")
	}

	err = mem.JoinGame(uuid.Nil.String(), uuid.NewString())
	if err == nil {
		t.Errorf("Expected error when joining nil game id")
	}

	minPlayers := 0
	maxPlayers := 8
	game, err := mem.CreateGame(minPlayers, maxPlayers)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = mem.JoinGame(game.ID.String(), "fake")
	if err == nil {
		t.Errorf("Expected error when unknown player joining game")
	}

	err = mem.JoinGame(game.ID.String(), uuid.Nil.String())
	if err == nil {
		t.Errorf("Expected error when nil id player joining game")
	}

	player1, err := mem.RegisterPlayer("", "name1")
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}

	err = mem.JoinGame(game.ID.String(), player1.ID.String())
	if err != nil {
		t.Errorf("error while player name1 joining game: %s", err.Error())
	}

	err = mem.JoinGame(game.ID.String(), player1.ID.String())
	if err == nil {
		t.Errorf("Expected error when already joined player joining game")
	}
}

func TestMemoryGames(t *testing.T) {
	logger := zerolog.Nop()

	// concrete memory test storage implementation
	mem := memory.New(&logger)

	minPlayers := 0
	maxPlayers := 4
	game, err := mem.CreateGame(minPlayers, maxPlayers)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if game == nil {
		t.Fatal("Expected game to be created, got nil")
	}

	id1 := game.ID.String()

	minPlayers = 0
	maxPlayers = 8
	game, err = mem.CreateGame(minPlayers, maxPlayers)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if game == nil {
		t.Error("Expected game to be created, got nil")
	}

	id2 := game.ID.String()

	games := mem.ListGames()
	if len(games) != 2 {
		t.Errorf("Expected 2 games, got %d", len(games))
	}

	// Start the game1 successfully
	err = mem.StartGame(id1)
	if err != nil {
		t.Errorf("Unexpected error when starting the game: %v", err)
	}

	started, err := mem.IsGameStarted(id1)
	if err != nil {
		t.Errorf("Unexpected error when checkin if game1 is started: %v", err)
	}

	if !started {
		t.Error("Expected game to be started, but it's not")
	}

	// Attempt to start the same game again, should return an error
	err = mem.StartGame(id1)
	if err == nil {
		t.Error("Expected error when starting an already started game")
	}

	// Attempt to stop a game never started, should return an error
	err = mem.StopGame(id2)
	if err == nil {
		t.Error("Expected error when stoppig a not started game")
	}

	// Start the game2 successfully
	err = mem.StartGame(id2)
	if err != nil {
		t.Errorf("Unexpected error when starting the game: %v", err)
	}

	started, err = mem.IsGameStarted(id2)
	if err != nil {
		t.Errorf("Unexpected error when checkin if game2 is started: %v", err)
	}

	if !started {
		t.Error("Expected game to be started, but it's not")
	}

	// Attempt to stop a game actually started successfully
	err = mem.StopGame(id2)
	if err != nil {
		t.Errorf("Unexpected error when stopping game2: %v", err)
	}

	started, err = mem.IsGameStarted(id2)
	if err != nil {
		t.Errorf("Unexpected error when checkin if game2 is started: %v", err)
	}

	if started {
		t.Error("Expected game to be stopped, but it's not")
	}

	// Attempt to stop the same game again, should return an error
	err = mem.StopGame(id2)
	if err == nil {
		t.Error("Expected error when stopping an already stopped game")
	}
}

func TestMemoryJoinGames(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[main] %s", i) },
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	// concrete memory test storage implementation
	mem := memory.New(&logger)

	minPlayers := 2
	maxPlayers := 4
	game, err := mem.CreateGame(minPlayers, maxPlayers)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	player1, err := mem.RegisterPlayer("", "name1")
	if err != nil {
		t.Errorf("error while registering name1: %s", err.Error())
	}

	player2, err := mem.RegisterPlayer("", "name2")
	if err != nil {
		t.Errorf("error while registering name2: %s", err.Error())
	}

	player3, err := mem.RegisterPlayer("", "name3")
	if err != nil {
		t.Errorf("error while registering name3: %s", err.Error())
	}

	player4, err := mem.RegisterPlayer("", "name4")
	if err != nil {
		t.Errorf("error while registering name4: %s", err.Error())
	}

	player5, err := mem.RegisterPlayer("", "name5")
	if err != nil {
		t.Errorf("error while registering name5: %s", err.Error())
	}

	err = mem.JoinGame(game.ID.String(), player1.ID.String())
	if err != nil {
		t.Errorf("error while player name1 joining game: %s", err.Error())
	}

	// Start the game with error
	err = mem.StartGame(game.ID.String())
	if err == nil {
		t.Errorf("game not supposed to start with less than 2 players")
	}

	err = mem.JoinGame(game.ID.String(), player2.ID.String())
	if err != nil {
		t.Errorf("error while player name2 joining game: %s", err.Error())
	}

	// Start the game successfully
	err = mem.StartGame(game.ID.String())
	if err != nil {
		t.Errorf("Unexpected error when starting the game: %v", err)
	}

	err = mem.JoinGame(game.ID.String(), player3.ID.String())
	if err == nil {
		t.Errorf("error expected as game is already started")
	}

	game2, err := mem.CreateGame(minPlayers, maxPlayers)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = mem.JoinGame(game2.ID.String(), player1.ID.String())
	if err != nil {
		t.Errorf("error while player name1 joining game: %s", err.Error())
	}
	err = mem.JoinGame(game2.ID.String(), player2.ID.String())
	if err != nil {
		t.Errorf("error while player name2 joining game: %s", err.Error())
	}
	err = mem.JoinGame(game2.ID.String(), player3.ID.String())
	if err != nil {
		t.Errorf("error while player name3 joining game: %s", err.Error())
	}
	err = mem.JoinGame(game2.ID.String(), player4.ID.String())
	if err != nil {
		t.Errorf("error while player name4 joining game: %s", err.Error())
	}
	err = mem.JoinGame(game2.ID.String(), player5.ID.String())
	if err == nil {
		t.Errorf("expected error while player name5 joining game because max players already reached")
	}

}
