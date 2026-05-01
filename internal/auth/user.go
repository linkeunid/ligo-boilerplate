package auth

// User represents the authenticated user.
type User struct {
	ID    string
	Name  string
	Email string
	Role  string
}

// IsAdmin checks if user has admin role.
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// IsRegular checks if user has regular user role.
func (u *User) IsRegular() bool {
	return u.Role == "user" || u.Role == ""
}
