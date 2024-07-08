package main

import (
	"database/sql"
	"time"
)

type Store interface {
	// for users
	CreateUser(*User) (*User, error)
	CreateTask(t *Task) (*Task, error)
	GetTask(id string) (*Task, error)
	GetUserById(id string) (*User, error)
	CreateProject(p *Project) (*Project, error)
	GetProjectById(id string) (*Project, error)
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec(`INSERT INTO users (firstName, lastName, email, password) 
    VALUES (? ? ? ?)`, u.FirstName, u.LastName, u.Email, u.Password,
	)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id

	return u, nil
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec(`INSERT INTO tasks (name, status, projectId, assignedToID) 
    VALUES (?, ?, ?, ?)`, t.Name, t.Status, t.ProjectID, t.AssignedToID,
	)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id
	t.CreatedAt = time.Now()

	return t, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task

	err := s.db.QueryRow(`SELECT id, name, status, projectId, assignedToID, createdAt 
    FROM tasks WHERE id = ?`, id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectID, &t.AssignedToID)

	return &t, err
}

func (s *Storage) GetUserById(id string) (*User, error) {
	var u User

	err := s.db.QueryRow(`SELECT id, firstName, lastName, password, email, createdAt, 
    FROM users WHERE id = ? `, id).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Password, &u.Email, &u.CreatedAt)

	return &u, err
}

func (s *Storage) CreateProject(p *Project) (*Project, error) {
	row, err := s.db.Exec("INTSER INTO projects (name) VALUES (?)", p.Name)
	if err != nil {
		return nil, err
	}

	id, err := row.LastInsertId()
	if err != nil {
		return nil, err
	}

	p.ID = id
	p.CreatedAt = time.Now()

	return p, nil
}

func (s *Storage) GetProjectById(id string) (*Project, error) {
	var p Project

	err := s.db.QueryRow(`SELECT id, name, createdAt FROM projects where id = ?`,
		id).Scan(&p.ID, &p.Name, &p.CreatedAt)

	return &p, err
}
