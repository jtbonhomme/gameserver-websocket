package storage

import (
	"github.com/jtbonhomme/gameserver-websocket/internal/games"
)

// Games defines the interface for games storage.
type Games interface {
	ListGames() []*games.Game                 // ListGames returns all games.
	CreateGame(int, int) (*games.Game, error) // CreateGame instantiates a new game.
	StartGame(string) error                   // StartGame starts the game with a given ID.
	StopGame(string) error                    // StopGame stops the game with a given ID.
	IsGameStarted(string) (bool, error)       // IsGameStarted returns true is game with given ID is started.
	JoinGame(string, string) error            // JoinGame adds a player to a game.
	GameByID(string) (*games.Game, error)     // GameByID returns a game object from its ID.
}
