package services

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/manabie-com/togo/internal/storages"
	postgres "github.com/manabie-com/togo/internal/storages/postgres"
	utils "github.com/manabie-com/togo/internal/utils"
)

// ToDoService implement HTTP server
type ToDoService struct {
	JWTKey string
	Store  *postgres.Pg
}

func (s *ToDoService) GetAuthToken(resp http.ResponseWriter, req *http.Request) {
	id := utils.Value(req, "user_id")
	if !s.Store.ValidateUser(req.Context(), id, utils.Value(req, "password")) {
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

func (s *ToDoService) createToken(id string) (string, error) {
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

func (s *ToDoService) ListTasks(resp http.ResponseWriter, req *http.Request) {
	id, _ := utils.UserIDFromCtx(req.Context())
	tasks, err := s.Store.RetrieveTasks(
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

func (s *ToDoService) AddTask(resp http.ResponseWriter, req *http.Request) {
	t := &storages.Task{}
	err := json.NewDecoder(req.Body).Decode(t)
	defer req.Body.Close()
	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	now := time.Now()
	userID, _ := utils.UserIDFromCtx(req.Context())
	t.ID = uuid.New().String()
	t.UserID = userID
	t.CreatedDate = now.Format("2006-01-02")

	// check current count of daily tasks
	count, err := s.Store.CountDailyTasks(req.Context(), sql.NullString{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	maxTasks, err := s.Store.GetMaxTasks(req.Context(), sql.NullString{
		String: userID,
		Valid:  true,
	})
	if err != nil {
		maxTasks = 5
	}

	if count > maxTasks {
		utils.RespondWithError(resp, http.StatusInternalServerError, "Exceeded max tasks per day")
		return
	}

	err = s.Store.AddTask(req.Context(), t)
	if err != nil {
		utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
		return
	}

	result := map[string]*storages.Task{"data": t}
	utils.RespondWithJSON(resp, http.StatusOK, result)
}




