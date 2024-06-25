package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("project id is required")
var errUserIDRequired = errors.New("user id is required")

type TasksService struct {
	store Store
}

func NewTasksService(s Store) *TasksService {
	return &TasksService{store: s}
}

func (s *TasksService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", WithJWTAuth(s.handleCreateTask, s.store)).Methods("POST")
	r.HandleFunc("/tasks/${id}", WithJWTAuth(s.handleGetTask, s.store)).Methods("GET")
}

func (s *TasksService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Failed to read request body"})
		return
	}
	defer r.Body.Close()

	// Declare a variable of type Task
	var task Task

	// Unmarshal JSON payload into the task variable
	if err := json.Unmarshal(body, &task); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid JSON format"})
		return
	}

	// Validate task payload
	if err := validateTaskPayload(&task); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Create task in the store
	createdTask, err := s.store.CreateTask(&task)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to create task"})
		return
	}

	// Respond with the created task
	WriteJSON(w, http.StatusCreated, createdTask)
}

func (s *TasksService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// if id == "" {
	// 	WriteJSON(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
	// 	return
	// }

	t, err := s.store.GetTask(id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "task not found!"})
		return
	}

	WriteJSON(w, http.StatusOK, t)
}

func validateTaskPayload(task *Task) error {
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
