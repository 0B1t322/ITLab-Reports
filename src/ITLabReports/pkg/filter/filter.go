package filter

import "golang.org/x/exp/constraints"

type FilterOperation int

func (f FilterOperation) GetOperationFilter() FilterOperation {
	return f
}

const (
	EQ FilterOperation = iota
	NEQ
	GT
	GTE
	LT
	LTE
	EXIST
	LIKE
)

type OperationFilterer interface {
	GetOperationFilter() FilterOperation
}

type FieldType interface {
	constraints.Float | constraints.Integer | ~string
}

type FilterField[T FieldType] struct {
	Value T
	Operation OperationFilterer
}

func (f FilterField[T]) GetValue() T {
	return f.Value
}

func (f FilterField[T]) GetOperation() FilterOperation {
	return f.Operation.GetOperationFilter()
}

type FieldFilterer[T FieldType] interface {
	GetValue() T
	GetOperation() FilterOperation
	BuildTo(builder FilterBuilder[T])
}

func (f *FilterField[T]) BuildTo(builder FilterBuilder[T]) {
	builder.BuildFilterField(f)
}

// FilterBuilder interface for exceute this filter on custom objects
type FilterBuilder[T FieldType] interface {
	BuildFilterField(
		field FieldFilterer[T],
	)
}