package constraints

import "golang.org/x/exp/constraints"

type BasicsTypes interface {
	constraints.Float | constraints.Integer | ~string
}