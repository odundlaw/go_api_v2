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
	ms := NewMockStore()
	router := mux.NewRouter()
	NewTasksRoutes(ms, MockJWTAuth).RegisterRoutes(router)

	t.Run("should return an error if name is empty", func(t *testing.T) {
		payload := &Task{
			Name: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer valid-token")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Error("invalid status code, it should fail")
		}
	})

	t.Run("should create a task", func(t *testing.T) {
		payload := &Task{
			Name:         "create task",
			AssignedToID: 45,
			ProjectID:    4,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer valid-token")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetTask(t *testing.T) {
	ms := NewMockStore()
	service := NewTasksRoutes(ms, MockJWTAuth)

	router := mux.NewRouter()
	service.RegisterRoutes(router)

	t.Run("it should return task", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/tasks/4", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer valid-token")

		rr := httptest.NewRecorder()

		ms.CreateTask(&Task{ID: 4, Name: "Odundlaw"})

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected staus %d, got status %d", http.StatusOK, rr.Code)
		}
	})
}
