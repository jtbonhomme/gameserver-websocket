package websocket

import (
	"time"

	"github.com/goombaio/namegenerator"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/pubsub"

	"github.com/jtbonhomme/test-gameserver-websocket/internal/manager"
)

type Server struct {
	log           *zerolog.Logger
	ps            *pubsub.Broker
	e             *echo.Echo
	nameGenerator namegenerator.Generator
	manager       *manager.Manager
}

func New(l *zerolog.Logger, m *manager.Manager) *Server {
	seed := time.Now().UTC().UnixNano()

	return &Server{
		log:           l,
		nameGenerator: namegenerator.NewNameGenerator(seed),
		manager:       m,
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

	s.manager.Start()
	s.e.Logger.Fatal(s.e.Start(":12345"))
}

func (s *Server) Shutdown() {
	s.e.Logger.Info("Server shuting down ...")
	// Shutdown
	s.manager.Shutdown()
}
