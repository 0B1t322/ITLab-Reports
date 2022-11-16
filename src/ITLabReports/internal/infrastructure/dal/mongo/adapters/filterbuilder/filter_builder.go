package filterbuilder

import (
	"fmt"

	"github.com/0B1t322/MongoBuilder/operators/aggregation"
	"github.com/0B1t322/MongoBuilder/operators/options"
	"github.com/0B1t322/MongoBuilder/operators/query"
	"github.com/0B1t322/MongoBuilder/utils"
	"github.com/0B1t322/RepoGen/pkg/filter"
	"go.mongodb.org/mongo-driver/bson"
)

type OperationBuilder interface {
	BuildFilterField(fieldName string, op filter.FilterOperation, fieldValue interface{}) bson.M
}

type AggragationOperationBuilder struct{}

func (a *AggragationOperationBuilder) BuildFilterField(
	fieldName string,
	op filter.FilterOperation,
	fieldValue interface{},
) bson.M {
	// Because aggragation work with $expr operator and fieldName should start with $
	fieldName = "$" + fieldName
	switch op {
	case filter.EQ:
		return aggregation.EQ(fieldName, fieldValue)
	case filter.GT:
		return aggregation.GT(fieldName, fieldValue)
	case filter.GTE:
		return aggregation.GTE(fieldName, fieldValue)
	case filter.LIKE:
		str := fmt.Sprint(fieldValue)
		return aggregation.RegexMatch(fieldName, str, options.I)
	case filter.LT:
		return aggregation.LT(fieldName, fieldValue)
	case filter.LTE:
		return aggregation.LTE(fieldName, fieldValue)
	case filter.NEQ:
		return aggregation.NE(fieldName, fieldValue)
	case filter.IN:
		return aggregation.In(fieldName, fieldValue)
	case filter.NIN:
		return aggregation.Not(aggregation.In(fieldName, fieldValue))
	}

	return bson.M{}
}

type FindOpeartionBuilder struct{}

func (f *FindOpeartionBuilder) BuildFilterField(
	fieldName string,
	op filter.FilterOperation,
	fieldValue interface{},
) bson.M {
	switch op {
	case filter.EQ:
		return query.EQ(fieldName, fieldValue)
	case filter.GT:
		return query.GT(fieldName, fieldValue)
	case filter.GTE:
		return query.GTE(fieldName, fieldValue)
	case filter.LIKE:
		str := fmt.Sprint(fieldValue)
		return query.Regex(fieldName, str, options.I)
	case filter.LT:
		return query.LT(fieldName, fieldValue)
	case filter.LTE:
		return query.LTE(fieldName, fieldValue)
	case filter.NEQ:
		return query.NE(fieldName, fieldValue)
	case filter.IN:
		if array, ok := fieldValue.(bson.A); ok {
			return query.In(fieldName, array...)
		}
		return query.In(fieldName, fieldValue)
	case filter.NIN:
		return query.Nin(fieldName, fieldValue)
	}

	return bson.M{}
}

type BuilderFilterAdapter[T filter.FieldType] struct {
	src              *bson.M
	Field            string
	FieldMarshaller  FieldMarshaller[T]
	OperationBuilder OperationBuilder
}

type FieldMarshaller[T filter.FieldType] func(field filter.FieldFilterer[T]) any

type BuilderFilterAdapterOptions[T filter.FieldType] func(b *BuilderFilterAdapter[T])

func WithFieldFormatter[T filter.FieldType](
	fieldMarshaller FieldMarshaller[T],
) BuilderFilterAdapterOptions[T] {
	return func(b *BuilderFilterAdapter[T]) {
		b.FieldMarshaller = fieldMarshaller
	}
}

func WithFindOperationBuilder[T filter.FieldType]() BuilderFilterAdapterOptions[T] {
	return func(b *BuilderFilterAdapter[T]) {
		b.OperationBuilder = &FindOpeartionBuilder{}
	}
}

// Set's default
func WithAggregationOperationBuilder[T filter.FieldType]() BuilderFilterAdapterOptions[T] {
	return func(b *BuilderFilterAdapter[T]) {
		b.OperationBuilder = &AggragationOperationBuilder{}
	}
}

func New[T filter.FieldType](
	to *bson.M,
	fieldName string,
	opts ...BuilderFilterAdapterOptions[T],
) filter.FilterBuilder[T] {
	b := &BuilderFilterAdapter[T]{
		src:              to,
		Field:            fieldName,
		OperationBuilder: &FindOpeartionBuilder{},
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func (b *BuilderFilterAdapter[T]) MarshallField(field filter.FieldFilterer[T]) any {
	if b.FieldMarshaller != nil {
		return b.FieldMarshaller(field)
	}

	return field.GetValue()
}

func (b *BuilderFilterAdapter[T]) BuildFilterField(
	field filter.FieldFilterer[T],
) {
	marshaledField := b.MarshallField(field)
	builded := b.OperationBuilder.BuildFilterField(b.Field, field.GetOperation(), marshaledField)
	*b.src = utils.MergeBsonM(*b.src, builded)
}
