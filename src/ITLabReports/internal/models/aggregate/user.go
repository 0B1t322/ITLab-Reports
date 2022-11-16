package aggregate

import "github.com/RTUITLab/ITLab-Reports/internal/models/entity"

type User struct {
	*entity.User
}

func NewUser(id string, role entity.UserRole) (User, error) {
	u := User{
		User: entity.NewUser(id, role),
	}

	return u, nil
}

func (u User) IsAdmin() bool {
	return u.Role == entity.UserRoleAdmin
}

func (u User) IsSuperAdmin() bool {
	return u.Role == entity.UserRoleSuperAdmin
}

func (u User) IsUser() bool {
	return u.Role == entity.UserRoleUser
}

func (u User) IsAdminOrSuperAdmin() bool {
	return u.IsAdmin() || u.IsSuperAdmin()
}
