package players

import (
	"github.com/jtbonhomme/pubsub"
	"github.com/rs/zerolog"
)

type Player struct {
	log    *zerolog.Logger
	client *pubsub.Client
}

func New(l *zerolog.Logger) *Player {
	return &Player{
		log: l,
	}
}
