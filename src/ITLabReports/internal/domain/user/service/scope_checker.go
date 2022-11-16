package user

import (
	"fmt"

	"github.com/samber/lo"
)

type ScopeChecker struct {
	RequiredScope string
}

func (sc *ScopeChecker) CheckScope(scopes []string) bool {
	_, find := lo.Find(
		scopes,
		func(v string) bool {
			return fmt.Sprint(v) == sc.RequiredScope
		},
	)

	return find
}
