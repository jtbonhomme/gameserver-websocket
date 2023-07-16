package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library

	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
)

const sqliteDatabaseFilepath string = "sqlite-database.db"

func main() {
	var err error
	// Init logger
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	output := zerolog.ConsoleWriter{
		Out:           os.Stderr,
		TimeFormat:    time.RFC3339,
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("[main] %s", i) },
	}
	logger := zerolog.New(output).With().Timestamp().Logger()

	_, err = os.Stat(sqliteDatabaseFilepath)
	if os.IsNotExist(err) {
		logger.Info().Msg("creating sqlite-database.db...")
		file, err := os.Create(sqliteDatabaseFilepath)
		if err != nil {
			logger.Panic().Msgf("error creating sqlite database: %w", err)
		}
		file.Close()
		logger.Info().Msg("sqlite-database.db created")
	} else {
		logger.Info().Msg("sqlite-database.db already exists")
	}
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer sqliteDatabase.Close()

	logger.Info().Msg("start manager")
	mgr := manager.New(&logger, sqliteDatabase)
	err = mgr.Start()
	if err != nil {
		logger.Panic().Msgf("manager start error: %w", err)
	}

	// todo: use websocket to send logs: provide an io.Writer implementing object
	// todo: manager websocket availability (because it starts in a goroutine)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info().Msg("received signal: " + s.String())
		// Shutdown
		logger.Info().Msg("start shutdown procedure")
		mgr.Shutdown()
	case err := <-mgr.Error():
		logger.Err(err).Msg("manager error")
	}
	logger.Info().Msg("exit")
}
