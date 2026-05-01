package user

import (
	"strconv"
	"sync"
)

// User represents a user entity.
type User struct {
	ID    string
	Name  string
	Email string
}

// UserRepository is the data access layer for users.
type UserRepository struct {
	mu     sync.RWMutex
	users  map[string]*User
	nextID int
}

// NewUserRepository creates a new user repository with sample data.
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[string]*User{
			"1": {ID: "1", Name: "Alice", Email: "alice@example.com"},
			"2": {ID: "2", Name: "Bob", Email: "bob@example.com"},
			"3": {ID: "3", Name: "Charlie", Email: "charlie@example.com"},
		},
		nextID: 4,
	}
}

// FindByID retrieves a user by ID.
func (r *UserRepository) FindByID(id string) (*User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, found := r.users[id]
	return user, found
}

// FindAll returns all users.
func (r *UserRepository) FindAll() []*User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

// Create adds a new user.
func (r *UserRepository) Create(name, email string) *User {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := strconv.Itoa(r.nextID)
	r.nextID++

	user := &User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	r.users[id] = user
	return user
}

// Update updates an existing user.
func (r *UserRepository) Update(id, name, email string) (*User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.users[id]; !found {
		return nil, false
	}

	user := &User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	r.users[id] = user
	return user, true
}

// Delete removes a user by ID.
func (r *UserRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.users[id]; !found {
		return false
	}

	delete(r.users, id)
	return true
}
