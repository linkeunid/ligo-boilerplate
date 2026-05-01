package todo

import (
	"sync"
)

// Repository handles data persistence for todos
type Repository struct {
	mu    sync.RWMutex
	todos map[string]*Todo
}

// NewRepository creates a new in-memory todo repository
func NewRepository() *Repository {
	return &Repository{
		todos: make(map[string]*Todo),
	}
}

// Save stores a todo
func (r *Repository) Save(todo *Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.todos[todo.ID] = todo
	return nil
}

// FindByID retrieves a todo by ID
func (r *Repository) FindByID(id string) (*Todo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	todo, ok := r.todos[id]
	return todo, ok
}

// FindAll returns all todos
func (r *Repository) FindAll() []*Todo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		result = append(result, todo)
	}
	return result
}

// Delete removes a todo by ID
func (r *Repository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.todos[id]; ok {
		delete(r.todos, id)
		return true
	}
	return false
}

// Exists checks if a todo ID exists
func (r *Repository) Exists(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.todos[id]
	return ok
}
