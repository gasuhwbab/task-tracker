package app

import (
	"log/slog"

	grpcapp "github.com/gasuhwbab/task-tracker/task-tracker/internal/app/grpc"
	tasktracker "github.com/gasuhwbab/task-tracker/task-tracker/internal/services/task-tracker"
	"github.com/gasuhwbab/task-tracker/task-tracker/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, gRPCPort int, storagePath string) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	taskTrackerService := tasktracker.New(log, storage, storage, storage, storage)
	grpcApp := grpcapp.New(log, taskTrackerService, gRPCPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
