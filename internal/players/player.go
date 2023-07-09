package players

import (
	"github.com/jtbonhomme/pubsub/client"
	"github.com/rs/zerolog"
)

type Player struct {
	log    *zerolog.Logger
	client *client.Client
}

func New(l *zerolog.Logger) *Player {
	return &Player{
		log: l,
	}
}
