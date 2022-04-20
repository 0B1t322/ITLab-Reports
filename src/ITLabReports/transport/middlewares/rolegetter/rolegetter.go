package rolegetter

import (
	"fmt"

	"golang.org/x/exp/slices"
)

// RoleGetter get roles from cliams
type RoleGetter struct {
	rolesWithPriority []string
}

func New(
	rolesWithPriority ...string,
) *RoleGetter {
	return &RoleGetter{
		rolesWithPriority: rolesWithPriority,
	}
}

func (r *RoleGetter) GetRole(claims map[string]any, claim string) (string, error) {
	var roles []string
	{
		switch claim := claims[claim].(type) {
		case []interface{}:
			for _, role := range claim {
				roles = append(roles, fmt.Sprint(role))
			}
		case interface{}:
			roles = append(roles, fmt.Sprint(claim))
		}
	}

	if len(roles) == 0 {
		return "", fmt.Errorf("Role not found")
	}

	slices.Sort(roles)

	var findedRole string = ""
	{
		for _, role := range r.rolesWithPriority {
			_, find := slices.BinarySearch(roles, role)
			if find {
				findedRole = role
				break
			}
		}
	}

	if findedRole == "" {
		return "", fmt.Errorf("Role not found")
	}

	return findedRole, nil

}


