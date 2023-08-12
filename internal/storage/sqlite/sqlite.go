package sqlite

import (
	"database/sql"
	"fmt"
)

type SQLite struct {
	db *sql.DB
}

// New creates a new SQLite object.
func New(db *sql.DB) (*SQLite, error) {
	s := &SQLite{
		db: db,
	}

	// Migrate sqlite tables
	err := s.MigrateSchema()
	if err != nil {
		return nil, fmt.Errorf("error migrating schema: %s", err.Error())
	}

	return s, nil
}
