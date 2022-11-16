package user_test

import (
	"testing"

	user "github.com/RTUITLab/ITLab-Reports/internal/domain/user/service"
	"github.com/stretchr/testify/require"
)

func TestFunc_ScopeCheckerTests(t *testing.T) {
	scopeChecker := user.ScopeChecker{RequiredScope: "this.events"}

	t.Run(
		"Find",
		func(t *testing.T) {
			require.True(
				t,
				scopeChecker.CheckScope(
					[]string{
						"this.events",
						"that.events",
						"some.service",
					},
				),
			)
		},
	)

	t.Run(
		"NotFind",
		func(t *testing.T) {
			require.False(
				t,
				scopeChecker.CheckScope(
					[]string{
						"that.events",
						"some.service",
					},
				),
			)
		},
	)
}
