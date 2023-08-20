package sqlite

import (
	"fmt"
)

// MigrateSchema migrates the database schema to the latest version.
func (s *SQLite) MigrateSchema() error {
	// Perform schema migration here, using a migration library like "goose" or "migrate"
	// This code snippet just demonstrates the concept, but you should use a proper migration tool in a real application.
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			uid UUID PRIMARY KEY,
			name VARCHAR(255) UNIQUE,
			score INT
		);

		CREATE TABLE IF NOT EXISTS games (
			uid UUID PRIMARY KEY,
			min_players INT,
			max_players INT,
			started BOOL,
			ended BOOL,
			start_time TIMESTAMP,
			end_time TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS game_players (
			game_id UUID,
			player_id UUID,
			FOREIGN KEY (game_id) REFERENCES games(id),
			FOREIGN KEY (player_id) REFERENCES players(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %s", err.Error())
	}

	return nil
}
