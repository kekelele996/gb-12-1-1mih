package constants

// UserRole enumerates platform roles.
const (
	RoleStudent   = "student"
	RoleCounselor = "counselor"
	RoleAdmin     = "admin"
)

// ValidRoles returns all accepted roles.
func ValidRoles() []string {
	return []string{RoleStudent, RoleCounselor, RoleAdmin}
}
