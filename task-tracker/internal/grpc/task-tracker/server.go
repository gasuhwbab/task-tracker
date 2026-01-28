package tasktracker

import (
	"context"

	tasktrackerv1 "github.com/gasuhwbab/task-tracker/protos/gen/go/task-tracker/v1"
)

type TaskTracker interface {
	CreateTask(context.Context)
	UpdateTask(context.Context)
	DeleteTask(context.Context)
	GetTasks(context.Context)
}

type serverAPI struct {
	tasktrackerv1.UnimplementedTaskTrackerServer
}

func (s *serverAPI) CreateTask() {

}

func (s *serverAPI) UpdatedTask() {

}

func (s *serverAPI) DeleteTask() {

}

func (s *serverAPI) GetTasks() {

}
