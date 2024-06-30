package main

import (
	"errors"
	"fmt"
	"net/http"
)

type MockStore struct {
	task map[string]*Task
	user map[string]*User
}

func NewMockStore() *MockStore {
	return &MockStore{
		task: make(map[string]*Task),
		user: make(map[string]*User),
	}
}

func (m *MockStore) CreateUser(u *User) (*User, error) {
	m.user[fmt.Sprintf("%d", u.ID)] = u
	return u, nil
}

func (m *MockStore) CreateTask(t *Task) (*Task, error) {
	m.task[fmt.Sprintf("%d", t.ID)] = t
	return t, nil
}

func (m *MockStore) GetUserById(id string) (*User, error) {
	if user, exists := m.user[id]; exists {
		return user, nil
	}
	return nil, errors.New("User does not exist")
}

func (m *MockStore) GetTask(id string) (*Task, error) {
	if task, exists := m.task[id]; exists {
		return task, nil
	}

	return nil, errors.New("Task does not exist")
}

func MockJWTAuth(handler http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer valid-token" {
			WriteJson(w, http.StatusForbidden, ErrorResponse{Error: "Unable to verify token"})
			return
		}

		handler.ServeHTTP(w, r)
	}
}
