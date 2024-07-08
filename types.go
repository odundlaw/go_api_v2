package main

import (
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type Task struct {
	CreatedAt    time.Time `json:"createdAt"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	ProjectID    int64     `json:"ProjectId"`
	AssignedToID int64     `json:"asssignedToId"`
	ID           int64     `json:"id"`
}

type User struct {
	CreatedAt time.Time `json:"createdAt"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	ID        int64     `json:"id"`
}

type Project struct {
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	ID        int64     `json:"id"`
}
