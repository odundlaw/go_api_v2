package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	errNameRequired      = errors.New("name is required")
	errProjectIDRequired = errors.New("project id is required")
	errUserIDRequired    = errors.New("user id is required")
)

type TaskService struct {
	store          Store
	authMiddleware func(http.HandlerFunc, Store) http.HandlerFunc
}

func NewTasksRoutes(store Store, middleware func(http.HandlerFunc, Store) http.HandlerFunc) *TaskService {
	return &TaskService{
		store:          store,
		authMiddleware: middleware,
	}
}

func (s *TaskService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", s.authMiddleware(s.handleCreateTask, s.store)).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{id}", s.authMiddleware(s.handleGetTask, s.store)).Methods(http.MethodGet)
}

func (s *TaskService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid requst Payload!"})
		return
	}

	defer r.Body.Close()

	var task *Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid requst Payload!"})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	task, err = s.store.CreateTask(task)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating Task"})
		return
	}

	WriteJson(w, http.StatusCreated, task)
}

func (s *TaskService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	id := GetParams(r, "id")

	if id == "" {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "id is required"})
		return
	}

	t, err := s.store.GetTask(id)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, http.StatusOK, t)
}

func validateTaskPayload(task *Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.AssignedToID == 0 {
		return errUserIDRequired
	}

	return nil
}
