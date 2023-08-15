package memory

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jtbonhomme/gameserver-websocket/internal/models"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

// ListGames returns all games.
func (m *Memory) ListGames() []*models.Game {
	games := []*models.Game{}
	for _, value := range m.games {
		games = append(games, value)
	}
	return games
}

// CreateGame instantiates a new game.
func (m *Memory) CreateGame(min, max int) (*models.Game, error) {
	gameID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to create game ID: %v", err)
	}

	game := &models.Game{
		ID:         gameID,
		MinPlayers: min,
		MaxPlayers: max,
		Started:    false,
		Players:    []string{},
	}
	m.games[gameID.String()] = game

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

	if game.Started {
		return fmt.Errorf("game already started")
	}

	players := game.Players
	if game.MinPlayers != 0 && len(players) < game.MinPlayers {
		return fmt.Errorf("min player number %d not reached yet", game.MinPlayers)
	}

	game.Started = true
	game.StartTime = time.Now()
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

	if !game.Started {
		return fmt.Errorf("game not started")
	}

	game.Started = false
	game.EndTime = time.Now()
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

	return game.Started, nil
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

	if game.Started {
		return fmt.Errorf("game already started, players can not join anymore")
	}

	players := game.Players
	if utils.ContainsString(players, idPlayer) {
		return fmt.Errorf("player id %s already joined the game", idPlayer)
	}

	if game.MaxPlayers != 0 && len(players) == game.MaxPlayers {
		return fmt.Errorf("max player number %d already already reached", game.MaxPlayers)
	}

	players = append(players, idPlayer)
	m.games[idGame].Players = players

	return nil
}
