package main

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(cfg mysql.Config) *MySQLStorage {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Db connection successful")

	return &MySQLStorage{
		db: db,
	}
}

func (s *MySQLStorage) Init() (*sql.DB, error) {
	// initialize the db tables
	if err := s.createUserTable(); err != nil {
		return nil, err
	}

	if err := s.createTasksTable(); err != nil {
		return nil, err
	}

	if err := s.createProjectTable(); err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *MySQLStorage) createProjectTable() error {
	_, err := s.db.Exec(`
    CREATE TABLE IF NOT EXISTS projects (
      id INT UNSIGNED NOT NULL AUTO_INCREMENT,
      name VARCHAR(255) NOT NULL,
      createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

      PRIMARY KEY (id)
    )
    ENGINE = InnoDB DEFAULT CHARACTER SET = utf8;
  `)

	return err
}

func (s *MySQLStorage) createUserTable() error {
	_, err := s.db.Exec(`
    CREATE TABLE IF NOT EXISTS users (
      id INT UNSIGNED NOT NULL AUTO_INCREMENT,
      email VARCHAR(255) NOT NULL,
      firstName VARCHAR(255) NOT NULL,
      lastName VARCHAR(255) NOT NULL,
      password VARCHAR(255) NOT NULL,
      createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
      
      PRIMARY KEY (id),
      UNIQUE INDEX email_UNIQUE (email ASC) VISIBLE
    )
    ENGINE = InnoDB DEFAULT CHARACTER SET = utf8
  `)

	return err
}

func (s *MySQLStorage) createTasksTable() error {
	_, err := s.db.Exec(`
    CREATE TABLE IF NOT EXISTS tasks  (
      id INT UNSIGNED NOT NULL AUTO_INCREMENT,
      name VARCHAR(255) NOT NULL,
      status ENUM('TODO', 'IN_PROGRESS', 'IN_TESTING', 'DONE') NOT NULL DEFAULT 'TODO',
      projectId INT UNSIGNED NOT NULL,
      assignedToID INT UNSIGNED NOT NULL,
      createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  
      PRIMARY KEY (id),
      INDEX projectId_idx (projectId ASC) VISIBLE,
      INDEX assignedToID_idx (assignedToID ASC) VISIBLE,
      
      CONSTRAINT projectId
      FOREIGN KEY (projectId) REFERENCES projects (id)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
  
      CONSTRAINT assignedToID
      FOREIGN KEY (assignedToID) REFERENCES users (id)
      ON DELETE CASCADE
      ON UPDATE CASCADE
    )
    ENGINE = InnoDB
    DEFAULT CHARACTER SET = utf8;
  `)

	return err
}
