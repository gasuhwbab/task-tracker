package main

import (
	"log/slog"

	"github.com/gasuhwbab/task-tracker/task-tracker/internal/config"
	"github.com/gasuhwbab/task-tracker/task-tracker/internal/logger"
)

func main() {
	// Init config
	cfg := config.MustLoad()
	// Init logger
	log := logger.New(cfg.Env)
	log.Info("starting app", slog.String("env", cfg.Env))
	
	// Init app

	// Run gRPC server
}
