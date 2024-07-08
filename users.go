package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var (
	errEmailRequired     = errors.New("valid email is required")
	errFirstNameRequired = errors.New("first name is required")
	errLastNameRequired  = errors.New("last name is required")
	errPasswordRequired  = errors.New("password is required")
)

type UserService struct {
	store Store
}

func NewUserService(store Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/create", s.handleCreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", s.handleGetUser).Methods(http.MethodGet)
}

func (s *UserService) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid requst Payload!"})
		return
	}

	defer r.Body.Close()

	payload := &User{}
	err = json.Unmarshal(body, payload)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid requst Payload!"})
		return
	}

	if err := validateUserPayload(payload); err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	hasedPassword, err := HashPassword(payload.Password)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	payload.Password = hasedPassword

	user, err := s.store.CreateUser(payload)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	WriteJson(w, http.StatusCreated, user)
}

func (s *UserService) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := GetParams(r, "id")
	if id == "" {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Insufficient request Params"})
		return
	}

	user, err := s.store.GetUserById(id)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, http.StatusOK, user)
}

func CreateAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := []byte(Envs.JWTSecret)

	token, err := CreateJWT(secret, id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})
	return token, nil
}

func validateUserPayload(user *User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if !strings.Contains(user.Email, "@") {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}
