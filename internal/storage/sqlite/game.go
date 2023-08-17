package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jtbonhomme/gameserver-websocket/internal/players"
)

// Game represents a game.
type Game struct {
	ID         int
	MinPlayers int
	MaxPlayers int
	Players    []*players.Player
	Started    bool
	StartTime  time.Time
	EndTime    time.Time
}

// CreateGame creates a new game with the specified minimum and maximum number of players.
func (s *SQLite) CreateGame(minPlayers, maxPlayers int) (*Game, error) {
	result, err := s.db.Exec("INSERT INTO games (min_players, max_players) VALUES (?, ?)", minPlayers, maxPlayers)
	if err != nil {
		return nil, fmt.Errorf("failed to create game: %v", err)
	}

	gameID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get game ID: %v", err)
	}

	game := &Game{
		ID:         int(gameID),
		MinPlayers: minPlayers,
		MaxPlayers: maxPlayers,
	}

	return game, nil
}

// JoinGame allows a player to join a game.
func (s *SQLite) JoinGame(gameID, playerID int) error {
	_, err := s.db.Exec("INSERT INTO game_players (game_id, player_id) VALUES (?, ?)", gameID, playerID)
	if err != nil {
		return fmt.Errorf("failed to join game: %v", err)
	}

	return nil
}

// StartGame starts a game when the required number of players is reached.
func (s *SQLite) StartGame(gameID int) error {
	_, err := s.db.Exec("UPDATE games SET started = true, start_time = NOW() WHERE id = ?", gameID)
	if err != nil {
		return fmt.Errorf("failed to start game: %v", err)
	}

	return nil
}

// RecordScore records the score for a player in a game.
func (s *SQLite) RecordScore(gameID, playerID, score int) error {
	_, err := s.db.Exec("INSERT INTO scores (game_id, player_id, score) VALUES (?, ?, ?)", gameID, playerID, score)
	if err != nil {
		return fmt.Errorf("failed to record score: %v", err)
	}

	return nil
}

// GameStats represents the statistics of a game.
type GameStats struct {
	Duration time.Duration
	Score    int
	// Add other statistics as needed
}

// GetGameStats retrieves the statistics for a game.
func (s *SQLite) GetGameStats(gameID int) (*GameStats, error) {
	var start time.Time
	var end time.Time
	var score int

	err := s.db.QueryRow("SELECT start_time, end_time, SUM(score) FROM games JOIN scores ON games.id = scores.game_id WHERE games.id = ? GROUP BY games.id", gameID).Scan(&start, &end, &score)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("game not found")
		}
		return nil, fmt.Errorf("failed to get game stats: %v", err)
	}

	duration := end.Sub(start)
	stats := &GameStats{
		Duration: duration,
		Score:    score,
	}

	return stats, nil
}

// HallOfFame represents the hall of fame of the best players.
type HallOfFame []*players.Player

// GetHallOfFame retrieves the hall of fame of the best players.
func (s *SQLite) GetHallOfFame(limit int) (HallOfFame, error) {
	rows, err := s.db.Query("SELECT players.id, players.name, SUM(scores.score) AS total_score FROM players JOIN scores ON players.id = scores.player_id GROUP BY players.id ORDER BY total_score DESC LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get hall of fame: %v", err)
	}
	defer rows.Close()

	var hallOfFame HallOfFame

	for rows.Next() {
		player := &players.Player{}
		err := rows.Scan(&player.ID, &player.Name, &player.Score)
		if err != nil {
			return nil, fmt.Errorf("failed to scan hall of fame row: %v", err)
		}

		hallOfFame = append(hallOfFame, player)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating hall of fame rows: %v", err)
	}

	return hallOfFame, nil
}
