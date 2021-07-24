package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	storages "github.com/trannguyenhung011086/togo/internal/storages"
	postgres "github.com/trannguyenhung011086/togo/internal/storages/postgres"
)

type TaskRepo struct {
	Store  *postgres.Pg
}

func (r *TaskRepo) ValidateUser(ctx context.Context, userID, pwd sql.NullString) bool {
 	res := r.Store.ValidateUser(ctx, userID, pwd)
	return res
}

func (r *TaskRepo) AddTask(ctx context.Context, task *storages.Task) (*storages.Task, error) {
	task.ID = uuid.New().String()
	task.CreatedDate = time.Now().Format("2006-01-02")
	
	err := r.Store.AddTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepo) ListTasks(ctx context.Context, userID, createdDate sql.NullString) ([]*storages.Task, error) {
	tasks, err := r.Store.RetrieveTasks(ctx, userID, createdDate)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func (r *TaskRepo) GetMaxTasks(ctx context.Context, userID sql.NullString) (int, error) {
	max, err := r.Store.GetMaxTasks(ctx, userID)
	if err != nil {
		return 0, err
	}
	return max, nil
}

func (r *TaskRepo) GetCurrentTasks(ctx context.Context, userID sql.NullString) (int, error) {
	count, err := r.Store.CountDailyTasks(ctx, userID)
	if err != nil {
		return 0, err
	}
	return count, nil
}