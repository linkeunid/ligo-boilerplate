package todo

import (
	"fmt"
	"time"
)

// Service handles business logic for todos
type Service struct {
	repo  *Repository
	nextID int
}

// NewService creates a new todo service
func NewService(repo *Repository) *Service {
	return &Service{
		repo:   repo,
		nextID: 1,
	}
}

// Create creates a new todo
func (s *Service) Create(req CreateTodoRequest) (*Todo, error) {
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	priority := req.Priority
	if priority == "" {
		priority = PriorityMedium
	}

	now := time.Now()
	todo := &Todo{
		ID:          fmt.Sprintf("%d", s.nextID),
		Title:       req.Title,
		Description: req.Description,
		Priority:    priority,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.nextID++

	if err := s.repo.Save(todo); err != nil {
		return nil, err
	}

	return todo, nil
}

// GetByID retrieves a todo by ID
func (s *Service) GetByID(id string) (*Todo, error) {
	todo, ok := s.repo.FindByID(id)
	if !ok {
		return nil, fmt.Errorf("todo not found")
	}
	return todo, nil
}

// List returns all todos
func (s *Service) List() []*Todo {
	return s.repo.FindAll()
}

// Update updates an existing todo
func (s *Service) Update(id string, req UpdateTodoRequest) (*Todo, error) {
	todo, ok := s.repo.FindByID(id)
	if !ok {
		return nil, fmt.Errorf("todo not found")
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	todo.UpdatedAt = time.Now()

	s.repo.Save(todo)
	return todo, nil
}

// Toggle toggles the completion status of a todo
func (s *Service) Toggle(id string) (*Todo, error) {
	todo, ok := s.repo.FindByID(id)
	if !ok {
		return nil, fmt.Errorf("todo not found")
	}

	todo.Completed = !todo.Completed
	todo.UpdatedAt = time.Now()

	s.repo.Save(todo)
	return todo, nil
}

// Delete deletes a todo by ID
func (s *Service) Delete(id string) error {
	if !s.repo.Delete(id) {
		return fmt.Errorf("todo not found")
	}
	return nil
}
