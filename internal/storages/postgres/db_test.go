package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	storages "github.com/trannguyenhung011086/togo/internal/storages"
)

var user = &storages.User{
	ID:       "firstUser",
	Password: "example",
}

var task = &storages.Task{
	ID: uuid.New().String(),
	Content: "test content",
	UserID: user.ID,
	CreatedDate: time.Now().Format("2006-01-02"),
}

func newMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestValidateUser(t *testing.T) {
	t.Run("Validate user failed", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		query := "SELECT id FROM users WHERE id = \\$1 AND password = \\$2"
		mock.ExpectQuery(query).WithArgs(user.ID, user.Password).WillReturnError(fmt.Errorf("Some error"))

		pg := &Pg{DB: db}
		ctx := context.Background()
		userId := sql.NullString{
			String: user.ID,
			Valid:  true,
		}
		pwd := sql.NullString{
			String: user.Password,
			Valid:  true,
		}

		res := pg.ValidateUser(ctx, userId, pwd)

		if res {
			t.Errorf("Expected validate user failed")
		}
	})

	t.Run("Validate user passed", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		query := "SELECT id FROM users WHERE id = \\$1 AND password = \\$2"
		rows := sqlmock.NewRows([]string{"id"}).AddRow(user.ID)
		mock.ExpectQuery(query).WithArgs(user.ID, user.Password).WillReturnRows(rows)

		pg := &Pg{DB: db}
		ctx := context.Background()
		userId := sql.NullString{
			String: user.ID,
			Valid:  true,
		}
		pwd := sql.NullString{
			String: user.Password,
			Valid:  true,
		}

		res := pg.ValidateUser(ctx, userId, pwd)

		if !res {
			t.Errorf("Expected validate user passed")
		}
	})
}

func TestAddTask(t *testing.T) {
	t.Run("Add task failed", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		query := "INSERT INTO tasks \\(id, content, user_id, created_date\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)"
		mock.ExpectExec(query).WithArgs(task.ID, task.Content, task.UserID, task.CreatedDate).WillReturnError(fmt.Errorf("Some error"))

		pg := &Pg{DB: db}
		ctx := context.Background()

		err := pg.AddTask(ctx, task)

		if err == nil {
			t.Errorf("Expected add task failed")
		}
	})

	t.Run("Add task passed", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		query := "INSERT INTO tasks \\(id, content, user_id, created_date\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\)"
		mock.ExpectExec(query).WithArgs(task.ID, task.Content, task.UserID, task.CreatedDate).WillReturnResult(sqlmock.NewResult(1, 1))

		pg := &Pg{DB: db}
		ctx := context.Background()

		err := pg.AddTask(ctx, task)

		if err != nil {
			t.Errorf("Expected add task passed")
		}
	})
}

func TestRetrieveTasks(t *testing.T) {
	t.Run("Get tasks failed", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		query := "SELECT id, content, user_id, created_date FROM tasks WHERE user_id = \\$1 AND created_date = \\$2"
		mock.ExpectQuery(query).WithArgs(user.ID, task.CreatedDate).WillReturnError(fmt.Errorf("Some error"))

		pg := &Pg{DB: db}
		ctx := context.Background()
		userId := sql.NullString{
			String: user.ID,
			Valid:  true,
		}
		createdDate := sql.NullString{
			String: task.CreatedDate,
			Valid:  true,
		}

		_, err := pg.RetrieveTasks(ctx, userId, createdDate)

		if err == nil {
			t.Errorf("Expected get tasks failed")
		}
	})

	t.Run("Get tasks passed", func(t *testing.T) {
		db, mock := newMock()
		defer db.Close()

		query := "SELECT id, content, user_id, created_date FROM tasks WHERE user_id = \\$1 AND created_date = \\$2"
		rows := sqlmock.NewRows([]string{"id", "content", "user_id", "created_date"}).AddRow(task.ID, task.Content, task.UserID, task.CreatedDate)
		mock.ExpectQuery(query).WithArgs(user.ID, task.CreatedDate).WillReturnRows(rows)

		pg := &Pg{DB: db}
		ctx := context.Background()
		userId := sql.NullString{
			String: user.ID,
			Valid:  true,
		}
		createdDate := sql.NullString{
			String: task.CreatedDate,
			Valid:  true,
		}

		_, err := pg.RetrieveTasks(ctx, userId, createdDate)

		if err != nil {
			t.Errorf("Expected get tasks passed")
		}
	})
}