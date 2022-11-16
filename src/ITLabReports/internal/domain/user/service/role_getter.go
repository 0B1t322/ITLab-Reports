package user

import (
	"github.com/samber/lo"
)

type OnRole struct {
	OnRole     string
	ReturnRole string
}

type RoleGetter struct {
	OnRoles []OnRole
	Default string
}

func (rg *RoleGetter) GetRole(roles []string) string {
	for _, onRole := range rg.OnRoles {
		if _, find := lo.Find(
			roles,
			func(v string) bool {
				return v == onRole.OnRole
			},
		); find {
			return onRole.ReturnRole
		}
	}

	return rg.Default
}
