package tasks

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AriJaya07/go-rest-api/packages/controllers/auth"
	"github.com/AriJaya07/go-rest-api/packages/models/store"
	"github.com/AriJaya07/go-rest-api/packages/models/types"
	"github.com/AriJaya07/go-rest-api/packages/utils"
	"github.com/gorilla/mux"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TasksService struct {
	store store.Store
}

func NewTasksService(s store.Store) *TasksService {
	return &TasksService{store: s}
}

func (s *TasksService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", auth.WithJWTAuth(s.handleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/${id}", auth.WithJWTAuth(s.handleGetTask, s.store)).Methods("GET")
}

func (s *TasksService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Declare a variable of type Task
	var task types.Task

	// Unmarshal JSON payload into the task variable
	if err := json.Unmarshal(body, &task); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	// Validate task payload
	if err := validateTaskPayload(&task); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
		return
	}

	// Create task in the store
	createdTask, err := s.store.CreateTask(&task)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create task"})
		return
	}

	// Respond with the created task
	utils.WriteJSON(w, http.StatusCreated, createdTask)
}

func (s *TasksService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// if id == "" {
	// 	utils.WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	// 	return
	// }

	t, err := s.store.GetTask(id)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, types.ErrorResponse{Error: "task not found!"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(task *types.Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectId == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}
