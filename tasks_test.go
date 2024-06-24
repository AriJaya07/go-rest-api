package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateTask(t *testing.T) {
	ms := &MockStore{}
	service := NewTasksService(ms)

	t.Run("should create a task", func(t *testing.T) {
		payload := &Task{
			Name:         "Creating a REST API in go",
			ProjectID:    1,
			AssignedToID: 42,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.handleCreateTask).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetTask(t *testing.T) {
	ms := &MockStore{}
	service := NewTasksService(ms)

	t.Run("should return the task", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/tasks/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks/{id}", service.handleGetTask)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Error("Invalid status code, it should fail")
		}
	})
}
