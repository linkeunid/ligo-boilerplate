package user

import (
	"regexp"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/common"
)

// CreateUserInput represents the input for creating a user.
type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validate checks if the input is valid.
func (i *CreateUserInput) Validate() error {
	if i.Name == "" {
		return common.ErrValidation
	}
	if len(i.Name) < 2 || len(i.Name) > 100 {
		return common.ErrValidation
	}
	if i.Email == "" {
		return common.ErrValidation
	}
	if !isValidEmail(i.Email) {
		return common.ErrValidation
	}
	return nil
}

// UpdateUserInput represents the input for updating a user.
type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Validate checks if the update input is valid.
func (i *UpdateUserInput) Validate() error {
	if i.Name != "" {
		if len(i.Name) < 2 || len(i.Name) > 100 {
			return common.ErrValidation
		}
	}
	if i.Email != "" {
		if !isValidEmail(i.Email) {
			return common.ErrValidation
		}
	}
	return nil
}

// UserService handles business logic for users.
type UserService struct {
	repo *UserRepository
	log  ligo.Logger
}

// NewUserService creates a new user service.
func NewUserService(repo *UserRepository, log ligo.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

// GetUserByID retrieves a user by ID.
func (s *UserService) GetUserByID(id string) (*User, error) {
	user, found := s.repo.FindByID(id)
	if !found {
		return nil, common.ErrNotFound
	}
	return user, nil
}

// GetAllUsers returns all users.
func (s *UserService) GetAllUsers() []*User {
	return s.repo.FindAll()
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(input CreateUserInput) (*User, error) {
	if err := input.Validate(); err != nil {
		s.log.Debug("User validation failed",
			ligo.LoggerField{Key: "error", Value: err.Error()},
		)
		return nil, err
	}

	user := s.repo.Create(input.Name, input.Email)
	s.log.Info("User created",
		ligo.LoggerField{Key: "user_id", Value: user.ID},
		ligo.LoggerField{Key: "name", Value: user.Name},
	)

	return user, nil
}

// UpdateUser updates an existing user.
func (s *UserService) UpdateUser(id string, input UpdateUserInput) (*User, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	name := input.Name
	email := input.Email

	// If fields are empty, keep existing values
	if name == "" || email == "" {
		if existing, found := s.repo.FindByID(id); found {
			if name == "" {
				name = existing.Name
			}
			if email == "" {
				email = existing.Email
			}
		} else {
			return nil, common.ErrNotFound
		}
	}

	user, updated := s.repo.Update(id, name, email)
	if !updated {
		return nil, common.ErrNotFound
	}

	s.log.Info("User updated",
		ligo.LoggerField{Key: "user_id", Value: user.ID},
	)

	return user, nil
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(id string) error {
	deleted := s.repo.Delete(id)
	if !deleted {
		return common.ErrNotFound
	}

	s.log.Info("User deleted",
		ligo.LoggerField{Key: "user_id", Value: id},
	)

	return nil
}

// isValidEmail checks if the email format is valid.
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
