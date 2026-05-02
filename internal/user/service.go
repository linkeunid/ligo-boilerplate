package user

import (
	"github.com/go-playground/validator/v10"
	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/common"
)

var validate = validator.New()

type CreateUserInput struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserInput struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UserService struct {
	repo *UserRepository
	log  ligo.Logger
}

func NewUserService(repo *UserRepository, log ligo.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

func (s *UserService) GetUserByID(id string) (*User, error) {
	user, found := s.repo.FindByID(id)
	if !found {
		return nil, common.ErrNotFound
	}
	return user, nil
}

func (s *UserService) GetAllUsers() []*User {
	return s.repo.FindAll()
}

func (s *UserService) CreateUser(input CreateUserInput) (*User, error) {
	if err := validate.Struct(input); err != nil {
		return nil, common.ErrValidation
	}

	user := s.repo.Create(input.Name, input.Email)
	s.log.Info("User created",
		ligo.LoggerField{Key: "user_id", Value: user.ID},
		ligo.LoggerField{Key: "name", Value: user.Name},
	)
	return user, nil
}

func (s *UserService) UpdateUser(id string, input UpdateUserInput) (*User, error) {
	if err := validate.Struct(input); err != nil {
		return nil, common.ErrValidation
	}

	name := input.Name
	email := input.Email

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
