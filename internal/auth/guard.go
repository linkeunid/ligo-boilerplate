package auth

import (
	"strings"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo-boilerplate/internal/common"
)

// Context keys for request-scoped data
const (
	ContextKeyUser = "user"
)

// AuthService handles authentication and authorization.
type AuthService struct {
	log ligo.Logger
}

// NewAuthService creates a new auth service.
func NewAuthService(log ligo.Logger) *AuthService {
	return &AuthService{log: log}
}

// ValidateToken validates a bearer token and returns the user.
// In production, this would validate a JWT token against an auth provider.
func (s *AuthService) ValidateToken(token string) (*User, error) {
	// TODO: Implement proper JWT validation
	// For demo purposes, accept tokens like "Bearer admin:token" or "Bearer user:token"

	if token == "" {
		return nil, common.ErrUnauthorized
	}

	// Demo: parse token format "role:secret"
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return nil, common.ErrUnauthorized
	}

	role := parts[0]

	s.log.Debug("Token validated",
		ligo.LoggerField{Key: "role", Value: role},
	)

	// Demo users based on role in token
	switch role {
	case "admin":
		return &User{
			ID:    "1",
			Name:  "Admin User",
			Email: "admin@example.com",
			Role:  "admin",
		}, nil
	case "user":
		return &User{
			ID:    "2",
			Name:  "Regular User",
			Email: "user@example.com",
			Role:  "user",
		}, nil
	default:
		return nil, common.ErrUnauthorized
	}
}

// AuthGuard validates the authorization header and stores user in context.
func AuthGuard(auth *AuthService) ligo.Guard {
	return func(ctx ligo.Context) (bool, error) {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return false, common.ErrUnauthorized
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return false, common.ErrUnauthorized
		}

		// Extract token without "Bearer " prefix
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		user, err := auth.ValidateToken(token)
		if err != nil {
			return false, common.ErrUnauthorized
		}

		// Store user in context for downstream handlers
		ctx.Set(ContextKeyUser, user)

		return true, nil
	}
}

// AdminGuard checks if user has admin role.
func AdminGuard() ligo.Guard {
	return func(ctx ligo.Context) (bool, error) {
		user, ok := ctx.Get(ContextKeyUser).(*User)
		if !ok || user == nil {
			return false, common.ErrUnauthorized
		}

		if !user.IsAdmin() {
			return false, common.ErrForbidden
		}

		return true, nil
	}
}
