package common

import "errors"

// Common error types for the application
var (
	// ErrUnauthorized is returned when authentication fails
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when authorization fails (valid user, insufficient permissions)
	ErrForbidden = errors.New("forbidden")

	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("not found")

	// ErrValidation is returned when input validation fails
	ErrValidation = errors.New("validation failed")
)
