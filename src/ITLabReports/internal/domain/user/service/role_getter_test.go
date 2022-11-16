package user_test

import (
	"testing"

	usersrv "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/stretchr/testify/require"
)

func TestFunc_RoleGetter(t *testing.T) {
	const (
		unknown    = "unknown"
		superadmin = "superadmin"
		admin      = "admin"
		user       = "user"
	)

	rg := usersrv.RoleGetter{
		OnRoles: []usersrv.OnRole{
			{
				OnRole:     "admin",
				ReturnRole: superadmin,
			},
			{
				OnRole:     "projects.admin",
				ReturnRole: admin,
			},
			{
				OnRole:     "projects.user",
				ReturnRole: user,
			},
		},
		Default: unknown,
	}

	t.Run(
		"FindSuperAdmin",
		func(t *testing.T) {
			require.Equal(
				t,
				superadmin,
				rg.GetRole(
					[]string{
						"projects.user",
						"service.role",
						"admin",
						"projects.admin",
					},
				),
			)
		},
	)

	t.Run(
		"FindAdmin",
		func(t *testing.T) {
			require.Equal(
				t,
				admin,
				rg.GetRole(
					[]string{
						"projects.user",
						"service.role",
						"projects.admin",
					},
				),
			)
		},
	)

	t.Run(
		"FindUser",
		func(t *testing.T) {
			require.Equal(
				t,
				user,
				rg.GetRole(
					[]string{
						"projects.user",
						"service.role",
					},
				),
			)
		},
	)

	t.Run(
		"Default",
		func(t *testing.T) {
			require.Equal(
				t,
				unknown,
				rg.GetRole(
					[]string{
						"service.role",
						"service.role",
						"service.role",
						"service.role",
					},
				),
			)
		},
	)
}
