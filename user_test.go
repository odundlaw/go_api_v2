package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestCreateUser(t *testing.T) {
	ms := NewMockStore()
	service := NewUserService(ms)

	router := mux.NewRouter()
	service.RegisterRoutes(router)

	t.Run("it should return an error if payload is invalid", func(t *testing.T) {
		payload := &User{
			Email: "",
		}

		bb, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/users/create", bytes.NewBuffer(bb))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected %d got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("it should return the created User", func(t *testing.T) {
		payload := &User{
			FirstName: "shittu",
			LastName:  "Lekan",
			Email:     "email@example.com",
			ID:        5,
			Password:  "dee22@",
			CreatedAt: time.Now(),
		}

		bb, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/users/create", bytes.NewBuffer(bb))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expceted %d got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetUser(t *testing.T) {
	ms := NewMockStore()
	router := mux.NewRouter()
	NewUserService(ms).RegisterRoutes(router)

	t.Run("it should return the user with the specified Id", func(t *testing.T) {
		payload := &User{
			FirstName: "shittu",
			LastName:  "Lekan",
			Email:     "email@example.com",
			ID:        5,
			Password:  "dee22@",
			CreatedAt: time.Now(),
		}

		_, err := ms.CreateUser(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodGet, "/users/5", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d go %d", http.StatusOK, rr.Code)
		}
	})
}
