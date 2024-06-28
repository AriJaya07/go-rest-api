package store

import (
	"database/sql"
	"log"

	"github.com/AriJaya07/go-rest-api/packages/models/types"
)

type Storage struct {
	db *sql.DB
}

type Store interface {
	// Users
	GetAllUsers() ([]types.User, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	CreateUser(u *types.User) (*types.User, error)
	GetUserByID(id string) (*types.User, error)
	UpdateUser(user *types.User) error
	UpdatePassword(user *types.User) error
	DeleteUser(id int64) error
	// Projects
	GetAllProjects() ([]types.Project, error)
	CreateProject(p *types.Project) error
	GetProject(id string) (*types.Project, error)
	UpdateProject(project *types.Project) error
	DeleteProject(id string) error
	// Tasks
	CreateTask(t *types.Task) (*types.Task, error)
	GetTask(id string) (*types.Task, error)
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) QueryRow(query string, args ...interface{}) *sql.Row {
	return s.db.QueryRow(query, args...)
}

func (s *Storage) GetAllUsers() ([]types.User, error) {
	var users []types.User

	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		log.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u types.User
		if err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Password, &u.CreatedAt); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error after scanning rows:", err)
		return nil, err
	}

	return users, nil
}

func (s *Storage) CreateUser(u *types.User) (*types.User, error) {
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

func (s *Storage) UpdateUser(user *types.User) error {
	// Prepare SQL statement
	query := "UPDATE users SET firstName = ?, lastName = ?, email = ? WHERE id = ?"

	// Execute SQL statement
	_, err := s.db.Exec(query, user.FirstName, user.LastName, user.Email, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetUserByID(id string) (*types.User, error) {
	var u types.User
	err := s.db.QueryRow("SELECT id, email, password, firstName, lastName, createdAt FROM users WHERE id = ?", id).Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.CreatedAt)
	return &u, err
}

func (s *Storage) DeleteUser(id int64) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func (s *Storage) CreateProject(p *types.Project) error {
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

func (s *Storage) GetAllProjects() ([]types.Project, error) {
	var projects []types.Project
	rows, err := s.db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p types.Project
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

func (s *Storage) UpdatePassword(user *types.User) error {
	query := "UPDATE users SET password = ? WHERE id = ?"
	_, err := s.db.Exec(query, user.Password, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetProject(id string) (*types.Project, error) {
	var p types.Project
	query := "SELECT id, name, createdAt FROM projects WHERE id = ?"
	err := s.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.CreatedAt)
	return &p, err
}

func (s *Storage) UpdateProject(project *types.Project) error {
	query := "UPDATE projects SET name = ? WHERE id = ?"

	_, err := s.db.Exec(query, project.Name, project.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteProject(id string) error {
	_, err := s.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) CreateTask(t *types.Task) (*types.Task, error) {
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

func (s *Storage) GetTask(id string) (*types.Task, error) {
	var t types.Task
	err := s.db.QueryRow("SELECT id, name, status, project_id, assigned_to, createdAt FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Name, &t.Status, &t.ProjectId, &t.AssignedToID, &t.CreatedAt)
	return &t, err
}
