package querybuilder

import (

	"golang.org/x/exp/constraints"
)

type QueryValue interface {
	constraints.Integer | constraints.Float | ~string
}

type MatchQuery[T QueryValue] struct {
	Field		string
	Value		T
}