package main

// Mocks

type MockStore struct{}

func (s *MockStore) CreateUser(u *User) (*User, error) {
	return &User{}, nil
}

func (s *MockStore) CreateTask(t *Task) (*Task, error) {
	return &Task{}, nil
}

func (s *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}
