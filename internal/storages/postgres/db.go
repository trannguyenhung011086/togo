package postgres

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	storages "github.com/trannguyenhung011086/togo/internal/storages"
)

type Pg struct {
	DB *sql.DB
}

func (p *Pg) ValidateUser(ctx context.Context, userID, pwd sql.NullString) bool {
	stmt := `SELECT id FROM users WHERE id = $1 AND password = $2`
	row := p.DB.QueryRowContext(ctx, stmt, userID, pwd)
	u := &storages.User{}
	err := row.Scan(&u.ID)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

func (p *Pg) AddTask(ctx context.Context, t *storages.Task) error {
	stmt := `INSERT INTO tasks (id, content, user_id, created_date) VALUES ($1, $2, $3, $4)`
	_, err := p.DB.ExecContext(ctx, stmt, &t.ID, &t.Content, &t.UserID, &t.CreatedDate)
	if err != nil {
		return err
	}

	return nil
} 

func (p *Pg) RetrieveTasks(ctx context.Context, userID, createdDate sql.NullString) ([]*storages.Task, error) {
	stmt := `SELECT id, content, user_id, created_date FROM tasks WHERE user_id = $1 AND created_date = $2`
	rows, err := p.DB.QueryContext(ctx, stmt, userID, createdDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*storages.Task
	for rows.Next() {
		t := &storages.Task{}
		err := rows.Scan(&t.ID, &t.Content, &t.UserID, &t.CreatedDate)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (p *Pg) GetMaxTasks(ctx context.Context, userID sql.NullString) (int, error) {
	stmt := `SELECT max_todo FROM users WHERE id = $1`
	row := p.DB.QueryRowContext(ctx, stmt, userID)

	var max int
	err := row.Scan(&max)
	if err != nil {
		return 0, err
	}
	return max, nil
}

func (p *Pg) CountDailyTasks(ctx context.Context, userID sql.NullString) (int, error) {
	currentDate:= time.Now().Format("2006-01-02")
	stmt := `SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND created_date = $2`
	rows, err := p.DB.QueryContext(ctx, stmt, userID, currentDate)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}