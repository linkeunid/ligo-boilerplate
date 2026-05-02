package memory

import (
	ligomemory "github.com/linkeunid/ligo-memory"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
)

// UserRepository is an in-memory implementation of repository.UserRepository
// backed by ligo-memory.Store.
type UserRepository struct {
	store *ligomemory.Store[int, *entity.User]
}

func NewUserRepository(store *ligomemory.Store[int, *entity.User]) repository.UserRepository {
	for _, u := range []*entity.User{
		{ID: 1, Name: "Alice", Email: "alice@example.com", Role: "user"},
		{ID: 2, Name: "Bob", Email: "bob@example.com", Role: "user"},
		{ID: 3, Name: "Charlie", Email: "charlie@example.com", Role: "admin"},
	} {
		store.Set(u.ID, u)
	}
	return &UserRepository{store: store}
}

func (r *UserRepository) FindByID(id int) (*entity.User, bool) {
	return r.store.Get(id)
}

func (r *UserRepository) FindAll() []*entity.User {
	return r.store.All()
}

func (r *UserRepository) Create(name, email, role string) *entity.User {
	id := nextID()
	user := &entity.User{ID: id, Name: name, Email: email, Role: role}
	r.store.Set(id, user)
	return user
}

func (r *UserRepository) Update(id int, name, email string) (*entity.User, bool) {
	existing, found := r.store.Get(id)
	if !found {
		return nil, false
	}
	updated := &entity.User{ID: id, Name: name, Email: email, Role: existing.Role}
	r.store.Set(id, updated)
	return updated, true
}

func (r *UserRepository) Delete(id int) bool {
	return r.store.Delete(id)
}
