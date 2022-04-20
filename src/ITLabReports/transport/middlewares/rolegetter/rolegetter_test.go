package rolegetter_test

import (
	"testing"

	"github.com/RTUITLab/ITLab-Reports/transport/middlewares/rolegetter"
	"github.com/stretchr/testify/require"
)


func TestFunc_RoleGetter(t *testing.T) {
	r := rolegetter.New(
		"admin",
		"service.admin",
		"service.user",
	)

	t.Run(
		"GetAdmin",
		func(t *testing.T) {
			role, err := r.GetRole(
				map[string]any{
					"roles": []any{
						"service.user",
						"service.admin",
						"admin",
					},
				},
				"roles",
			)

			require.NoError(t, err)
			require.Equal(t, "admin", role)
		},
	)

	t.Run(
		"GetServiceAdmin",
		func(t *testing.T) {
			role, err := r.GetRole(
				map[string]any{
					"roles": []any{
						"service.user",
						"service.admin",
					},
				},
				"roles",
			)

			require.NoError(t, err)
			require.Equal(t, "service.admin", role)
		},
	)

	t.Run(
		"GetServiceUser",
		func(t *testing.T) {
			role, err := r.GetRole(
				map[string]any{
					"roles": []any{
						"service.user",
					},
				},
				"roles",
			)

			require.NoError(t, err)
			require.Equal(t, "service.user", role)
		},
	)

	t.Run(
		"NotFoundrole",
		func(t *testing.T) {
			_, err := r.GetRole(
				map[string]any{
					"roles": []any{
						"some_role",
					},
				},
				"roles",
			)

			require.Error(t, err)
		},
	)
}