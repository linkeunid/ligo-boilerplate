package user

import (
	"github.com/linkeunid/ligo"
)

// CreateUserValidationPipe validates and transforms create user input.
func CreateUserValidationPipe(log ligo.Logger) ligo.Pipe {
	return func(input any) (any, error) {
		log.Debug("Validating create user input")
		// In real app, use proper validation (e.g., go-playground/validator)
		return input, nil
	}
}

// UpdateUserValidationPipe validates and transforms update user input.
func UpdateUserValidationPipe(log ligo.Logger) ligo.Pipe {
	return func(input any) (any, error) {
		log.Debug("Validating update user input")
		// In real app, use proper validation (e.g., go-playground/validator)
		return input, nil
	}
}
