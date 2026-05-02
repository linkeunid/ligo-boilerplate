package auth

// User represents the authenticated user.
type User struct {
	ID    string
	Name  string
	Email string
	Role  string
}

// HasRole implements the ligo.HasRole interface for role-based guards.
func (u *User) HasRole(role string) bool {
	return u.Role == role
}

// IsAdmin checks if user has admin role.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsRegular checks if user has regular user role.
func (u *User) IsRegular() bool {
	return u.Role == "user" || u.Role == ""
}
