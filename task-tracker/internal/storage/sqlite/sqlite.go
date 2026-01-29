package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	tasktrackerv1 "github.com/gasuhwbab/task-tracker/protos/gen/go/task-tracker/v1"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveUser saves user to db.
func (s *Storage) CreateTask(ctx context.Context, name string) (uint32, error) {
	const op = "storage.sqlite.SaveTask"

	stmt, err := s.db.Prepare("INSERT INTO tasks(name) VALUES(?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, name)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return uint32(id), nil
}

func (s *Storage) UpdateTask(ctx context.Context, taskId uint32, name string, text string, progress uint32) (uint32, error) {
	const op = "storage.sqlite.UpdateTask"

	stmt, err := s.db.Prepare("UPDATE tasks SET (name, text, progress) VALUES(?, ?, ?) WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if _, err := stmt.ExecContext(ctx, name, text, progress, taskId); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return taskId, nil
}

func (s *Storage) DeleteTask(ctx context.Context, taskId uint32) (uint32, error) {
	const op = "storage.sqlite.DeleteTask"

	stmt, err := s.db.Prepare("DELETE FROM tasks WHERE id = ?")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if _, err := stmt.ExecContext(ctx, taskId); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return taskId, nil
}

func (s *Storage) GetTasks(ctx context.Context) ([]*tasktrackerv1.Task, error) {
	const op = "storage.sqlite.GetTasks"

	stmt, err := s.db.Prepare("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	res := []*tasktrackerv1.Task{}
	for rows.Next() {
		var task tasktrackerv1.Task
		if err := rows.Scan(&task); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, &task)
	}
	return res, nil
}
