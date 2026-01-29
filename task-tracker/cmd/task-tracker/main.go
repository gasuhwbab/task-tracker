package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gasuhwbab/task-tracker/task-tracker/internal/app"
	"github.com/gasuhwbab/task-tracker/task-tracker/internal/config"
	"github.com/gasuhwbab/task-tracker/task-tracker/internal/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)
	log.Info("starting app", slog.String("env", cfg.Env))
	app := app.New(log, cfg.GRPC.Port, cfg.StoragePath)
	go app.GRPCServer.MustRun()
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	stopSignal := <-stopChan
	log.Info("stopping application", slog.String("stop signal", stopSignal.String()))
	app.GRPCServer.Stop()
}
