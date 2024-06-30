package main

import "database/sql"

type Store interface {
	// for users
	CreateUser(*User) (*User, error)
	CreateTask(t *Task) (*Task, error)
	GetTask(id string) (*Task, error)
	GetUserById(id string) (*User, error)
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
