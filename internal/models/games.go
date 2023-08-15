package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID         uuid.UUID `json:"id"`
	Started    bool      `json:"started"`
	Players    []string  `json:"players"`
	MinPlayers int       `json:"minPlayers"`
	MaxPlayers int       `json:"maxPlayers"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
}
