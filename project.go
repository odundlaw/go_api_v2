package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errName = errors.New("name cannot be empty")

type ProjectService struct {
	store Store
}

func NewProjectService(s Store) *ProjectService {
	return &ProjectService{store: s}
}

func (p *ProjectService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/projects/create", p.handleCreateProject).Methods(http.MethodPost)
	r.HandleFunc("/projects/{id}", p.handleGetProject).Methods(http.MethodPost)
}

func (p *ProjectService) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid Request Body"})
		return
	}

	defer r.Body.Close()

	var project *Project

	if err := json.Unmarshal(body, project); err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid Paylaod"})
		return
	}

	if err := validateProject(project); err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	pr, err := p.store.CreateProject(project)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJson(w, http.StatusCreated, pr)
}

func (p *ProjectService) handleGetProject(w http.ResponseWriter, r *http.Request) {
	id := GetParams(r, "id")

	if id == "" {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Insufficient Paramaters"})
		return
	}

	project, err := p.store.GetProjectById(id)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJson(w, http.StatusOK, project)
}

func validateProject(p *Project) error {
	if p.Name == "" {
		return errName
	}

	return nil
}
