package main

import (
	"database/sql"
)

type Storage struct {
	db *sql.DB
}

type Store interface {
	// Users
	CreateUser(u *User) (*User, error)
	GetUserByID(id string) (*User, error)
	// Projects
	GetAllProjects() ([]Project, error)
	CreateProject(p *Project) error
	GetProject(id string) (*Project, error)
	DeleteProject(id string) error
	// Tasks
	CreateTask(t *Task) (*Task, error)
	GetTask(id string) (*Task, error)
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateProject(p *Project) error {
	result, err := s.db.Exec("INSERT INTO projects (name) VALUES (?)", p.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = id

	return err
}

func (s *Storage) GetAllProjects() ([]Project, error) {
	var projects []Project
	rows, err := s.db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *Storage) GetProject(id string) (*Project, error) {
	var p Project
	query := "SELECT id, name, createdAt FROM projects WHERE id = ?"
	err := s.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.CreatedAt)
	return &p, err
}

func (s *Storage) DeleteProject(id string) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec("INSERT INTO users (email, firstName, lastName, password) VALUES (?, ?, ?, ?)", u.Email, u.FirstName, u.LastName, u.Password)
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

func (s *Storage) GetUserByID(id string) (*User, error) {
	var u User
	err := s.db.QueryRow("SELECT id, email, firstName, lastName, createdAt FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt)
	return &u, err
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec("INSERT INTO tasks (name, status, project_id, assigned_to) VALUES (?, ?, ?, ?)", t.Name, t.Status, t.ProjectId, t.AssignedToID)

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
	err := s.db.QueryRow("SELECT id, name, status, project_id, assigned_to, createdAt FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectId, &t.AssignedToID, &t.CreatedAt)
	return &t, err
}
