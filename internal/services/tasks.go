package services

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	repository "github.com/trannguyenhung011086/togo/internal/repository"
	"github.com/trannguyenhung011086/togo/internal/storages"
	utils "github.com/trannguyenhung011086/togo/internal/utils"
)


type TaskService struct {
	JWTKey string
	TaskRepo *repository.TaskRepo
}

func (s *TaskService) GetAuthToken(resp http.ResponseWriter, req *http.Request) {
	id := utils.Value(req, "user_id")
	if !s.TaskRepo.ValidateUser(req.Context(), id, utils.Value(req, "password")) {
		utils.RespondWithError(resp, http.StatusUnauthorized, "incorrect user_id/pwd")
		return
	}

	token, err := s.createToken(id.String)
	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]string{"data": token}
	utils.RespondWithJSON(resp, http.StatusOK, result)
}

func (s *TaskService) createToken(id string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(s.JWTKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *TaskService) ListTasks(resp http.ResponseWriter, req *http.Request) {
	id, _ := utils.UserIDFromCtx(req.Context())
	tasks, err := s.TaskRepo.ListTasks(
		req.Context(),
		sql.NullString{
			String: id,
			Valid:  true,
		},
		utils.Value(req, "created_date"),
	)

	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string][]*storages.Task{"data": tasks}
	utils.RespondWithJSON(resp, http.StatusOK, result)
}

func (s *TaskService) AddTask(resp http.ResponseWriter, req *http.Request) {
	// parse task data
	task := &storages.Task{}
	err := json.NewDecoder(req.Body).Decode(task)
	defer req.Body.Close()
	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	userID, ok := utils.UserIDFromCtx(req.Context())
	if !ok {
		utils.RespondWithError(resp, http.StatusInternalServerError, "Invalid userId")
		return
	}
	
	task.UserID = userID
	
	// check current count of daily tasks
	count, err := s.TaskRepo.GetCurrentTasks(req.Context(), sql.NullString{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	maxTasks, err := s.TaskRepo.GetMaxTasks(req.Context(), sql.NullString{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		maxTasks = 5
	}

	if count >= maxTasks {
		utils.RespondWithError(resp, http.StatusInternalServerError, "Exceeded max tasks per day")
		return
	}

	// add task
	task, err = s.TaskRepo.AddTask(req.Context(), task)
	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]*storages.Task{"data": task}
	utils.RespondWithJSON(resp, http.StatusOK, result)
}
