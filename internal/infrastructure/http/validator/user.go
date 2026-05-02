package validator

import (
	"github.com/linkeunid/ligo-boilerplate/internal/usecase/dto"
)

const maxFileSize = 10 << 20 // 10MB

// UserValidator handles request validation for user endpoints.
type UserValidator struct{}

// NewUserValidator creates a new user validator.
func NewUserValidator() *UserValidator {
	return &UserValidator{}
}

// ValidateCreateUser validates create user request.
func (v *UserValidator) ValidateCreateUser(input *dto.CreateUserInput) error {
	// Validation is now handled by use case layer with struct tags
	// This validator can be extended for custom validation logic
	return nil
}

// ValidateUpdateUser validates update user request.
func (v *UserValidator) ValidateUpdateUser(input *dto.UpdateUserInput) error {
	return nil
}

// Common validation errors.
var (
	ErrFileTooLarge = &ValidationError{Message: "file size exceeds 10MB limit"}
)

// ValidationError represents a validation error.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
