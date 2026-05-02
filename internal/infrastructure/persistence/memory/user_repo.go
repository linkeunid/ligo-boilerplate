package memory

import (
	"crypto/rand"
	"fmt"
	"sync"

	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

// UserRepository is an in-memory implementation of repository.UserRepository.
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*entity.User
}

// NewUserRepository creates a new in-memory user repository with sample data.
func NewUserRepository() repository.UserRepository {
	return &UserRepository{
		users: map[string]*entity.User{
			"550e8400-e29b-41d4-a716-446655440001": {
				ID: "550e8400-e29b-41d4-a716-446655440001", Name: "Alice", Email: "alice@example.com", Role: "user",
			},
			"550e8400-e29b-41d4-a716-446655440002": {
				ID: "550e8400-e29b-41d4-a716-446655440002", Name: "Bob", Email: "bob@example.com", Role: "user",
			},
			"550e8400-e29b-41d4-a716-446655440003": {
				ID: "550e8400-e29b-41d4-a716-446655440003", Name: "Charlie", Email: "charlie@example.com", Role: "admin",
			},
		},
	}
}

func (r *UserRepository) FindByID(id string) (*entity.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, found := r.users[id]
	return user, found
}

func (r *UserRepository) FindAll() []*entity.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users
}

func (r *UserRepository) Create(name, email string) *entity.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := newUUID()
	user := &entity.User{
		ID:    id,
		Name:  name,
		Email: email,
		Role:  "user",
	}

	r.users[id] = user
	return user
}

func (r *UserRepository) Update(id, name, email string) (*entity.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, found := r.users[id]
	if !found {
		return nil, false
	}

	updated := &entity.User{
		ID:    id,
		Name:  name,
		Email: email,
		Role:  existing.Role,
	}

	r.users[id] = updated
	return updated, true
}

func (r *UserRepository) Delete(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, found := r.users[id]; !found {
		return false
	}

	delete(r.users, id)
	return true
}

func newUUID() string {
	var b [16]byte
	rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
