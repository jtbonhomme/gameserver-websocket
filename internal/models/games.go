package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID                     uuid.UUID `json:"id"`
	Started                bool
	Players                []string
	MinPlayers, MaxPlayers int
	StartTime              time.Time
	EndTime                time.Time
}
