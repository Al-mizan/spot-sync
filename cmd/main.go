package main

import (
	"log/slog"
	"os"
	"spotsync/internal/config"
	"spotsync/internal/server"
)

func main() {
	// load environment variables
	cfg, err := config.LoadEnv()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// connect to the database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		slog.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	// start the server
	server.Start(db, cfg)
}
