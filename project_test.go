package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateProject(t *testing.T) {
	ms := NewMockStore()
	router := mux.NewRouter()
	NewProjectService(ms).RegisterRoutes(router)

	t.Run("it should return error for invalid payload", func(t *testing.T) {
		payload := &Project{
			Name: "",
		}

		bb, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/projects/create", bytes.NewBuffer(bb))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected error code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}
