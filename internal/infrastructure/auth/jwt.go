package auth

import (
	"strings"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/entity"
	"github.com/linkeunid/ligo-boilerplate/internal/domain/service"
	"github.com/linkeunid/ligo-boilerplate/internal/usecase"
)

// JWTAuth implements domain.service.AuthService.
// For demo purposes, accepts tokens in format "role:secret".
// Replace with actual JWT validation in production.
type JWTAuth struct {
	log ligo.Logger
}

// NewJWTAuth creates a new JWT auth service.
func NewJWTAuth(log ligo.Logger) *JWTAuth {
	return &JWTAuth{log: log}
}

// ValidateToken validates a bearer token and returns the authenticated user.
func (j *JWTAuth) ValidateToken(token string) (service.User, error) {
	if token == "" {
		return nil, usecase.ErrUnauthorized
	}

	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return nil, usecase.ErrUnauthorized
	}

	role := parts[0]
	j.log.Debug("Token validated", ligo.LoggerField{Key: "role", Value: role})

	switch role {
	case "admin":
		return &entity.User{
			ID:    3,
			Name:  "Admin User",
			Email: "admin@example.com",
			Role:  entity.RoleAdmin,
		}, nil
	case "user":
		return &entity.User{
			ID:    1,
			Name:  "Regular User",
			Email: "user@example.com",
			Role:  entity.RoleUser,
		}, nil
	default:
		return nil, usecase.ErrUnauthorized
	}
}

// OnModuleInit is called when the auth module initializes.
// Use this for validation checks, loading keys, etc.
func (j *JWTAuth) OnModuleInit() error {
	j.log.Info("JWT authentication initialized")
	// In production, validate JWT keys, load certificates, etc.
	return nil
}

// OnApplicationShutdown is called during application shutdown.
// Use this for cleanup, closing connections, etc.
func (j *JWTAuth) OnApplicationShutdown() error {
	j.log.Info("JWT authentication shutting down")
	// In production, close connections, release resources, etc.
	return nil
}
