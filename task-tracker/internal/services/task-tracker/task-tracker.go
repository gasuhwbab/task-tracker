package tasktracker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	tasktrackerv1 "github.com/gasuhwbab/task-tracker/protos/gen/go/task-tracker/v1"
)

var (
	ErrTaskExists   = errors.New("task already exists")
	ErrTaskNotExist = errors.New("task doesn't exists")
)

type TaskTracker struct {
	log         *slog.Logger
	taskCreator TaskCreator
	taskUpdater TaskUpdater
	taskGetter  TaskGetter
	taskDeleter TaskDeleter
}

type TaskCreator interface {
	CreateTask(context.Context, string) (uint32, error)
}

type TaskUpdater interface {
	UpdateTask(context.Context, uint32, string, string, uint32) (uint32, error)
}

type TaskDeleter interface {
	DeleteTask(context.Context, uint32) (uint32, error)
}

type TaskGetter interface {
	GetTasks(context.Context) ([]*tasktrackerv1.Task, error)
}

func New(
	log *slog.Logger,
	taskCreator TaskCreator,
	taskUpdater TaskUpdater,
	taskDeleter TaskDeleter,
	taskGetter TaskGetter,
) *TaskTracker {
	return &TaskTracker{
		log:         log,
		taskCreator: taskCreator,
		taskUpdater: taskUpdater,
		taskGetter:  taskGetter,
		taskDeleter: taskDeleter,
	}
}

func (taskTracker *TaskTracker) CreateTask(ctx context.Context, userId uint32, name string) (uint32, error) {
	const op = "tasktracker.CreateTask"
	log := taskTracker.log.With(
		slog.String("op", op),
		slog.Int64("user_id", int64(userId)),
	)
	log.Info("creating task")
	taskId, err := taskTracker.taskCreator.CreateTask(ctx, name)
	if err != nil {
		if errors.Is(err, ErrTaskExists) {
			log.Warn("task already exists", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		log.Warn("failed to create task", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return taskId, nil
}

func (taskTracker *TaskTracker) UpdateTask(
	ctx context.Context,
	userId uint32,
	newTask *tasktrackerv1.Task,
) (uint32, error) {
	const op = "tasktracker.UpdateTask"
	log := taskTracker.log.With(
		slog.String("op", op),
		slog.Int64("user_id", int64(userId)),
	)
	log.Info("updating task")
	taskId, err := taskTracker.taskUpdater.UpdateTask(ctx, newTask.TaskId, newTask.Name, newTask.Description, newTask.Progress)
	if err != nil {
		log.Warn("failed to update task", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return taskId, nil
}

func (taskTracker *TaskTracker) DeleteTask(
	ctx context.Context,
	userId uint32,
	taskId uint32,
) (uint32, error) {
	const op = "tasktracker.DeleteTask"
	log := taskTracker.log.With(
		slog.String("op", op),
		slog.Int64("user_id", int64(userId)),
	)
	log.Info("deleting task")
	taskId, err := taskTracker.taskDeleter.DeleteTask(ctx, taskId)
	if err != nil {
		if errors.Is(err, ErrTaskNotExist) {
			log.Warn("task doesn't exists", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		log.Warn("failed to delete task", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return taskId, nil
}

func (taskTracker *TaskTracker) GetTasks(
	ctx context.Context,
	userId uint32,
) ([]*tasktrackerv1.Task, error) {
	const op = "tasktracker.GetTasks"
	log := taskTracker.log.With(
		slog.String("op", op),
		slog.Int64("user_id", int64(userId)),
	)
	log.Info("getting tasks")
	tasks, err := taskTracker.taskGetter.GetTasks(ctx)
	if err != nil {
		log.Warn("failed to getting tasks", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}
