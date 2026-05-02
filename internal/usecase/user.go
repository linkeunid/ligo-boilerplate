package usecase

import (
	"github.com/go-playground/validator/v10"
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/repository"
	"github.com/linkeunid/ligo-boilerplate/internal/usecase/dto"
)

// UserUseCase contains business logic for user operations.
type UserUseCase struct {
	repo   repository.UserRepository
	log    ligo.Logger
	verify *validator.Validate
}

// NewUserUseCase creates a new user use case.
func NewUserUseCase(repo repository.UserRepository, log ligo.Logger) *UserUseCase {
	return &UserUseCase{
		repo:   repo,
		log:    log,
		verify: validator.New(),
	}
}

// GetUserByID retrieves a user by ID.
func (uc *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	user, found := uc.repo.FindByID(id)
	if !found {
		return nil, ErrNotFound
	}
	return user, nil
}

// GetAllUsers returns all users.
func (uc *UserUseCase) GetAllUsers() []*entity.User {
	return uc.repo.FindAll()
}

// CreateUser creates a new user with validation.
func (uc *UserUseCase) CreateUser(input dto.CreateUserInput) (*entity.User, error) {
	if err := uc.verify.Struct(input); err != nil {
		return nil, ErrValidation
	}

	user := uc.repo.Create(input.Name, input.Email)
	uc.log.Info("User created",
		ligo.LoggerField{Key: "user_id", Value: user.ID},
		ligo.LoggerField{Key: "name", Value: user.Name},
	)
	return user, nil
}

// UpdateUser updates an existing user with validation.
func (uc *UserUseCase) UpdateUser(id string, input dto.UpdateUserInput) (*entity.User, error) {
	if err := uc.verify.Struct(input); err != nil {
		return nil, ErrValidation
	}

	name := input.Name
	email := input.Email

	// Preserve existing values for partial updates
	if name == "" || email == "" {
		if existing, found := uc.repo.FindByID(id); found {
			if name == "" {
				name = existing.Name
			}
			if email == "" {
				email = existing.Email
			}
		} else {
			return nil, ErrNotFound
		}
	}

	user, updated := uc.repo.Update(id, name, email)
	if !updated {
		return nil, ErrNotFound
	}

	uc.log.Info("User updated",
		ligo.LoggerField{Key: "user_id", Value: user.ID},
	)
	return user, nil
}

// DeleteUser deletes a user by ID.
func (uc *UserUseCase) DeleteUser(id string) error {
	deleted := uc.repo.Delete(id)
	if !deleted {
		return ErrNotFound
	}
	uc.log.Info("User deleted",
		ligo.LoggerField{Key: "user_id", Value: id},
	)
	return nil
}
