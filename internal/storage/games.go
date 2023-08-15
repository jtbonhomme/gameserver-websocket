package storage

import (
	"github.com/jtbonhomme/gameserver-websocket/internal/models"
)

// Games defines the interface for games storage.
type Games interface {
	ListGames() []*models.Game                 // ListGames returns all games.
	CreateGame(int, int) (*models.Game, error) // CreateGame instantiates a new game.
	StartGame(string) error                    // StartGame starts the game with a given ID.
	StopGame(string) error                     // StopGame stops the game with a given ID.
	IsGameStarted(string) (bool, error)        // IsGameStarted returns true is game with given ID is started.
	JoinGame(string, string) error             // JoinGame joins a player to a game.
}
