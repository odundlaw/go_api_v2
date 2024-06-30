package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	store Store
	add   string
}

func NewApiServer(add string, store Store) *APIServer {
	return &APIServer{
		add:   add,
		store: store,
	}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	NewUserService(s.store).RegisterRoutes(subrouter)
	NewTasksRoutes(s.store, WithJJWTAuth).RegisterRoutes(subrouter)

	log.Print("Server listening at ", s.add)

	log.Fatal(http.ListenAndServe(s.add, subrouter))
}
