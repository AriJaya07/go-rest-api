package store_test

import (
	"errors"

	"github.com/AriJaya07/go-rest-api/packages/models/types"
)

// Mocks

type MockStore struct{}

func (s *MockStore) CreateProject(p *types.Project) error {
	return nil
}

func (ms *MockStore) GetAllProjects() ([]types.Project, error) {
	return nil, errors.New("Not implemented")
}

func (s *MockStore) GetProject(id string) (*types.Project, error) {
	return &types.Project{Name: "Super cool project"}, nil
}

func (s *MockStore) DeleteProject(id string) error {
	return nil
}

func (s *MockStore) CreateUser(u *types.User) (*types.User, error) {
	return &types.User{}, nil
}

func (s *MockStore) GetUserByID(id string) (*types.User, error) {
	return &types.User{}, nil
}

func (s *MockStore) CreateTask(t *types.Task) (*types.Task, error) {
	return &types.Task{}, nil
}

func (s *MockStore) GetTask(id string) (*types.Task, error) {
	return &types.Task{}, nil
}
