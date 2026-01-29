package tasktracker

import (
	"context"

	tasktrackerv1 "github.com/gasuhwbab/task-tracker/protos/gen/go/task-tracker/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyVal = 0
)

type TaskTracker interface {
	CreateTask(context.Context, uint32, string) (uint32, error)
	UpdateTask(context.Context, uint32, *tasktrackerv1.Task) (uint32, error)
	DeleteTask(context.Context, uint32, uint32) (uint32, error)
	GetTasks(context.Context, uint32) ([]*tasktrackerv1.Task, error)
}

type serverAPI struct {
	tasktrackerv1.UnimplementedTaskTrackerServer
	TaskTracker
}

func Register(gRPC *grpc.Server, taskTracker TaskTracker) {
	tasktrackerv1.RegisterTaskTrackerServer(gRPC, &serverAPI{TaskTracker: taskTracker})
}

func (s *serverAPI) CreateTask(
	ctx context.Context,
	req *tasktrackerv1.CreateTaskRequest,
) (*tasktrackerv1.CreateTaskResponse, error) {
	if err := validateCreateTask(req); err != nil {
		return nil, err
	}
	taskId, err := s.TaskTracker.CreateTask(ctx, req.GetUserId(), req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create task")
	}
	return &tasktrackerv1.CreateTaskResponse{TaskId: taskId}, nil
}

func (s *serverAPI) UpdateTask(
	ctx context.Context,
	req *tasktrackerv1.UpdateTaskRequest,
) (*tasktrackerv1.UpdateTaskResponse, error) {
	if err := validateUpdateTask(req); err != nil {
		return nil, err
	}
	taskId, err := s.TaskTracker.UpdateTask(ctx, req.GetUserId(), req.GetUpdatedTask())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update task")
	}
	return &tasktrackerv1.UpdateTaskResponse{TaskId: taskId}, nil
}

func (s *serverAPI) DeleteTask(
	ctx context.Context,
	req *tasktrackerv1.DeleteTaskRequest,
) (*tasktrackerv1.DeleteTaskResponse, error) {
	if err := validateDeleteTask(req); err != nil {
		return nil, err
	}
	taskId, err := s.TaskTracker.DeleteTask(ctx, req.GetUserId(), req.GetTaskId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete task")
	}
	return &tasktrackerv1.DeleteTaskResponse{TaskId: taskId}, nil
}

func (s *serverAPI) GetTasks(
	ctx context.Context,
	req *tasktrackerv1.GetTasksRequest,
) (*tasktrackerv1.GetTasksResponse, error) {
	if err := validageGetTasks(req); err != nil {
		return nil, err
	}
	tasks, err := s.TaskTracker.GetTasks(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get tasks")
	}
	return &tasktrackerv1.GetTasksResponse{Tasks: tasks}, nil
}

func validateCreateTask(req *tasktrackerv1.CreateTaskRequest) error {
	if req.GetUserId() == emptyVal {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	return nil
}

func validateUpdateTask(req *tasktrackerv1.UpdateTaskRequest) error {
	if req.GetUserId() == emptyVal {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetUpdatedTask() == nil {
		return status.Error(codes.InvalidArgument, "task is required")
	}
	return nil
}

func validateDeleteTask(req *tasktrackerv1.DeleteTaskRequest) error {
	if req.GetUserId() == emptyVal {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetTaskId() == emptyVal {
		return status.Error(codes.InvalidArgument, "task_id is required")
	}
	return nil
}

func validageGetTasks(req *tasktrackerv1.GetTasksRequest) error {
	if req.GetUserId() == emptyVal {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
