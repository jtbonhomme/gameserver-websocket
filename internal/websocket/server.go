package websocket

import (
	"time"

	"github.com/goombaio/namegenerator"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/pubsub"
)

type Server struct {
	log           *zerolog.Logger
	ps            *pubsub.Broker
	e             *echo.Echo
	nameGenerator namegenerator.Generator
}

func New(l *zerolog.Logger, ps *pubsub.Broker) *Server {
	seed := time.Now().UTC().UnixNano()

	return &Server{
		log:           l,
		ps:            ps,
		nameGenerator: namegenerator.NewNameGenerator(seed),
	}
}

func (s *Server) Start() {
	// 1. Server Setup
	s.log.Info().Msg("Server starting ...")

	s.ps = pubsub.New(s.log)

	s.e = echo.New()
	s.e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	s.e.Use(middleware.Recover())
	s.e.Static("/", "panel")
	s.e.GET("/connect", s.connect)

	s.e.Logger.Fatal(s.e.Start(":12345"))
}

func (s *Server) Shutdown() {
	s.e.Logger.Info("Server shuting down ...")
}
