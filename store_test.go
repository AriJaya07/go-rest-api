package main

import "errors"

// Mocks

type MockStore struct{}

func (s *MockStore) CreateProject(p *Project) error {
	return nil
}

func (ms *MockStore) GetAllProjects() ([]Project, error) {
	return nil, errors.New("Not implemented")
}

func (s *MockStore) GetProject(id string) (*Project, error) {
	return &Project{Name: "Super cool project"}, nil
}

func (s *MockStore) DeleteProject(id string) error {
	return nil
}

func (s *MockStore) CreateUser(u *User) (*User, error) {
	return &User{}, nil
}

func (s *MockStore) GetUserByID(id string) (*User, error) {
	return &User{}, nil
}

func (s *MockStore) CreateTask(t *Task) (*Task, error) {
	return &Task{}, nil
}

func (s *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}
