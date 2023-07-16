package manager

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jtbonhomme/gameserver-websocket/internal/game"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"

	"github.com/centrifugal/centrifuge"
	"github.com/rs/zerolog"
)

const defaultShutdownTimeout = 3 * time.Second

type Manager struct {
	log             *zerolog.Logger
	games           []*game.Game
	player          []players.Player
	started         bool
	err             chan error
	node            *centrifuge.Node
	shutdownTimeout time.Duration
	db              *sql.DB
}

func auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Put authentication Credentials into request Context.
		// Since we don't have any session backend here we simply
		// set user ID as empty string. Users with empty ID called
		// anonymous users, in real app you should decide whether
		// anonymous users allowed to connect to your server or not.
		cred := &centrifuge.Credentials{
			UserID: "",
		}
		newCtx := centrifuge.SetCredentials(ctx, cred)
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

// New creates a new Manager instance.
func New(l *zerolog.Logger, db *sql.DB) *Manager {
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[manager] %s", i) },
	}
	logger := l.Output(output)

	return &Manager{
		log:             &logger,
		games:           []*game.Game{},
		err:             make(chan error),
		shutdownTimeout: defaultShutdownTimeout,
		db:              db,
	}
}

// Error sends an error in the dedicated channel.
func (m *Manager) Error() <-chan error {
	return m.err
}

// Start starts the manager.
func (m *Manager) Start() error {
	m.log.Info().Msg("starting ...")

	node, err := centrifuge.New(centrifuge.Config{})
	if err != nil {
		return fmt.Errorf("error creating centrifuge node: %w", err)
	}
	m.node = node
	m.node.OnConnect(func(client *centrifuge.Client) {
		// In our example transport will always be Websocket but it can also be SockJS.
		transportName := client.Transport().Name()
		// In our example clients connect with JSON protocol but it can also be Protobuf.
		transportProto := client.Transport().Protocol()
		m.log.Info().Msgf("client %s (%s) connected via %s (%s)", client.ID(), string(client.Info()), transportName, transportProto)

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			m.log.Info().Msgf("client %s (%s) subscribes on channel %s", client.ID(), string(client.Info()), e.Channel)
			cb(centrifuge.SubscribeReply{}, nil)
		})

		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			m.log.Info().Msgf("client %s (%s) publishes into channel %s: %s", client.ID(), string(client.Info()), e.Channel, string(e.Data))
			cb(centrifuge.PublishReply{}, nil)
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			m.log.Info().Msgf("client %s (%s) disconnected", client.ID(), string(client.Info()))
		})

		client.OnRPC(m.HandleRPC)
	})

	// Run node. This method does not block. See also node.Shutdown method
	// to finish application gracefully.
	if err := m.node.Run(); err != nil {
		return fmt.Errorf("error running centrifuge node: %w", err)
	}

	// Configure HTTP routes.
	// Serve Websocket connections using WebsocketHandler.
	wsHandler := centrifuge.NewWebsocketHandler(m.node, centrifuge.WebsocketConfig{})
	http.Handle("/connection/websocket", auth(wsHandler))

	// The second route is for serving index.html file.
	http.Handle("/", http.FileServer(http.Dir("./public")))

	go func() {
		m.log.Info().Msgf("Starting server, visit http://localhost:8000")
		if err := http.ListenAndServe(":8000", nil); err != nil {
			m.err <- fmt.Errorf("error listening on :8000: %w", err)
		}
	}()

	// Migrate sqlite tables
	m.MigrateSchema()

	return nil
}

func (m *Manager) Shutdown() {
	m.log.Info().Msg("shuting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
	defer cancel()

	m.node.Shutdown(ctx)

	m.log.Info().Msgf("stopped")
}
