package entity

type UserRole string

func (u UserRole) String() string {
	return string(u)
}

const (
	UserRoleAdmin      UserRole = "admin"
	UserRoleUser       UserRole = "user"
	UserRoleSuperAdmin UserRole = "super.admin"
	UserRoleUnknown    UserRole = "unknown"
)

// UserRole from string
func UserRoleFromString(s string) UserRole {
	switch s {
	case "admin":
		return UserRoleAdmin
	case "super.admin":
		return UserRoleSuperAdmin
	case "user":
		return UserRoleUser
	default:
		return UserRoleUnknown
	}
}

type User struct {
	ID   string
	Role UserRole
}

func NewUser(id string, role UserRole) *User {
	return &User{
		ID:   id,
		Role: role,
	}
}
