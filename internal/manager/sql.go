package manager

import (
	"fmt"
)

// MigrateSchema migrates the database schema to the latest version.
func (m *Manager) MigrateSchema() error {
	// Perform schema migration here, using a migration library like "goose" or "migrate"
	// This code snippet just demonstrates the concept, but you should use a proper migration tool in a real application.
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255)
		);

		CREATE TABLE IF NOT EXISTS games (
			id INT AUTO_INCREMENT PRIMARY KEY,
			min_players INT,
			max_players INT,
			started BOOL,
			start_time TIMESTAMP,
			end_time TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS game_players (
			game_id INT,
			player_id INT,
			FOREIGN KEY (game_id) REFERENCES games(id),
			FOREIGN KEY (player_id) REFERENCES players(id)
		);

		CREATE TABLE IF NOT EXISTS scores (
			id INT AUTO_INCREMENT PRIMARY KEY,
			game_id INT,
			player_id INT,
			score INT,
			FOREIGN KEY (game_id) REFERENCES games(id),
			FOREIGN KEY (player_id) REFERENCES players(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	return nil
}
