package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	repository "github.com/trannguyenhung011086/togo/internal/repository"
	services "github.com/trannguyenhung011086/togo/internal/services"
	postgres "github.com/trannguyenhung011086/togo/internal/storages/postgres"
	utils "github.com/trannguyenhung011086/togo/internal/utils"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
	TaskService *services.TaskService
}

func (a *App) run(addr string) {
	log.Fatal(http.ListenAndServe(":5050", a.Router))
}

func (app *App) initialize(user, password, dbname, host, port, jwtkey string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", user, password, dbname, host, port)

	var err error
	app.DB, err = sql.Open("postgres", connectionString)
	if err == nil {
		log.Printf("success connecting db %s", connectionString)
	} else {
		log.Fatal("error connecting db", err)
	}

	stmt, err := app.DB.Prepare(`CREATE TABLE IF NOT EXISTS users (
		id TEXT NOT NULL,
		password TEXT NOT NULL,
		max_todo INTEGER DEFAULT 5 NOT NULL,
		CONSTRAINT users_PK PRIMARY KEY (id)
	);`)
	if err != nil {
		log.Fatal("error preparing db", err)
	}
	stmt.Exec()
	stmt, err = app.DB.Prepare(`CREATE TABLE IF NOT EXISTS tasks (
		id TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id TEXT NOT NULL,
		created_date TEXT NOT NULL,
		CONSTRAINT tasks_PK PRIMARY KEY (id),
		CONSTRAINT tasks_FK FOREIGN KEY (user_id) REFERENCES users(id)
	);`)
	if err != nil {
		log.Fatal("error preparing db", err)
	}
	stmt.Exec()
	defer stmt.Close()

	app.Router = mux.NewRouter()

    app.initializeRoutes(jwtkey)
}

func (app *App) initializeRoutes(jwtkey string) {
	TaskService := &services.TaskService{JWTKey: jwtkey, TaskRepo: &repository.TaskRepo{Store: &postgres.Pg{DB: app.DB}}}
	
	app.Router.HandleFunc("/login", TaskService.GetAuthToken).Methods("GET")
	app.Router.HandleFunc("/tasks", authMiddleware(jwtkey, TaskService.ListTasks)).Methods("GET")
	app.Router.HandleFunc("/tasks", authMiddleware(jwtkey, TaskService.AddTask)).Methods("POST")
}

func authMiddleware(jwtkey string, next http.HandlerFunc) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")

		claims := make(jwt.MapClaims)
		t, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
			return []byte(jwtkey), nil
		})
		if err != nil {
			utils.RespondWithError(resp, http.StatusUnauthorized, err.Error())
			return
		}

		if !t.Valid {
			utils.RespondWithError(resp, http.StatusUnauthorized, err.Error())
			return
		}

		id, ok := claims["user_id"].(string)
		if !ok {
			utils.RespondWithError(resp, http.StatusUnauthorized, err.Error())
			return
		}

		req = req.WithContext(context.WithValue(req.Context(), utils.UserAuthKey(0), id))
		next.ServeHTTP(resp, req)
	}
}

func main() {
	app := App{}

	app.initialize(os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_PORT"),
		os.Getenv("JWT_KEY"))

	app.run(":5050")
}
