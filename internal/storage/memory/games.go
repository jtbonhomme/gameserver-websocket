package memory

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/jtbonhomme/gameserver-websocket/internal/games"
)

// ListGames returns all games.
func (m *Memory) ListGames() []*games.Game {
	result := []*games.Game{}
	for _, value := range m.games {
		result = append(result, value)
	}
	return result
}

// CreateGame instantiates a new game.
func (m *Memory) CreateGame(min, max int) (*games.Game, error) {
	game := games.New(m.log, min, max)

	err := game.Connect()
	if err != nil {
		return nil, fmt.Errorf("error creating game: %s", err.Error())
	}
	m.games[game.ID.String()] = game

	return game, nil
}

// StartGame starts the game with a given ID.
func (m *Memory) StartGame(id string) error {
	if id == uuid.Nil.String() {
		return fmt.Errorf("nil game id")
	}

	game, ok := m.games[id]
	if !ok {
		return fmt.Errorf("unknown game id %s", id)
	}

	err := game.Start()
	if err != nil {
		return fmt.Errorf("error starting game: %s", err.Error())
	}
	m.games[id] = game

	return nil
}

// StopGame stops the game with a given ID.
func (m *Memory) StopGame(id string) error {
	if id == uuid.Nil.String() {
		return fmt.Errorf("nil game id")
	}

	game, ok := m.games[id]
	if !ok {
		return fmt.Errorf("unknown game id %s", id)
	}

	err := game.Stop()
	if err != nil {
		return fmt.Errorf("error stopping game: %s", err.Error())
	}

	m.games[id] = game

	return nil
}

// IsGameStarted returns true is game with given ID is started.
func (m *Memory) IsGameStarted(id string) (bool, error) {
	if id == uuid.Nil.String() {
		return false, fmt.Errorf("nil game id")
	}

	game, ok := m.games[id]
	if !ok {
		return false, fmt.Errorf("unknown game id %s", id)
	}

	return game.IsStarted(), nil
}

// JoinGame adds a player to a game.
func (m *Memory) JoinGame(idGame, idPlayer string) error {
	if idGame == uuid.Nil.String() {
		return fmt.Errorf("nil game id")
	}

	game, ok := m.games[idGame]
	if !ok {
		return fmt.Errorf("unknown game id %s", idGame)
	}

	if idPlayer == uuid.Nil.String() {
		return fmt.Errorf("nil player id")
	}

	_, ok = m.players[idPlayer]
	if !ok {
		return fmt.Errorf("unknown player id %s", idPlayer)
	}

	err := game.AddPlayer(idPlayer)
	if err != nil {
		return fmt.Errorf("error adding player to game: %s", err.Error())
	}
	m.games[idGame] = game

	return nil
}

// GameByID returns a game object from its ID.
func (m *Memory) GameByID(id string) (*games.Game, error) {
	// if provided id matches an existing game
	if id != uuid.Nil.String() {
		game, ok := m.games[id]
		if ok {
			return game, nil
		}
	}

	return nil, fmt.Errorf("unknown game id: %s", id)
}
