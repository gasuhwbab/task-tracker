package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	tasktrackerv1 "github.com/gasuhwbab/task-tracker/protos/gen/go/task-tracker/v1"
	tasktracker "github.com/gasuhwbab/task-tracker/task-tracker/internal/grpc/task-tracker"
	"google.golang.org/grpc"
)

type TaskTracker interface {
	CreateTask(context.Context, uint32, string) (uint32, error)
	UpdateTask(context.Context, uint32, *tasktrackerv1.Task) (uint32, error)
	DeleteTask(context.Context, uint32, uint32) (uint32, error)
	GetTasks(context.Context, uint32) ([]*tasktrackerv1.Task, error)
}

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, taskTrackerService TaskTracker, port int) *App {
	gRPCServer := grpc.NewServer()
	tasktracker.Register(gRPCServer, taskTrackerService)
	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().Network()))
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op), slog.Int("port", a.port))
	a.gRPCServer.GracefulStop()
}
