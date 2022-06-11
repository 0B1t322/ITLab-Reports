package mongobuildertofilter

import (
	"fmt"
	"time"

	"github.com/RTUITLab/ITLab-Reports/pkg/dialect/mongo"
	builder "github.com/RTUITLab/ITLab-Reports/pkg/dialect/mongo"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuilderFilterAdapter[T filter.FieldType] struct {
	Builder *builder.Predicate
	Field   string

	FieldMarshaller FieldMarshaller[T]
}

type FieldMarshaller[T filter.FieldType] func(field filter.FieldFilterer[T]) any

func FieldToTime[T filter.FieldType]() FieldMarshaller[T] {
	return func(field filter.FieldFilterer[T]) any {
		str := fmt.Sprint(field.GetValue())

		time, err := time.Parse(time.RFC3339, str)
		if err != nil {
			return nil
		}

		return time.UTC()
	}
}

func StringIdMarshaller() FieldMarshaller[string] {
	return func(field filter.FieldFilterer[string]) any {
		str := fmt.Sprint(field.GetValue())

		id, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			return nil
		}

		return id
	}
}

func SliceIdMarshaller() FieldMarshaller[[]string] {
	return func(field filter.FieldFilterer[[]string]) any {
		result := mongo.Array()
		for _, elem := range field.GetValue() {
			id, err := primitive.ObjectIDFromHex(elem)
			if err != nil {
				continue
			}
			result.AddElem(id)
		}
		return result.Array()
	}
}

type BuilderFilterAdapterOptions[T filter.FieldType] func(b *BuilderFilterAdapter[T])

func WithFieldFormatter[T filter.FieldType](fieldMarshaller FieldMarshaller[T]) BuilderFilterAdapterOptions[T] {
	return func(b *BuilderFilterAdapter[T]) {
		b.FieldMarshaller = fieldMarshaller
	}
}

func New[T filter.FieldType](
	builder *builder.Predicate,
	fieldName string,
	opts ...BuilderFilterAdapterOptions[T],
) filter.FilterBuilder[T] {
	b := &BuilderFilterAdapter[T]{
		Builder: builder,
		Field:   fieldName,
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

	switch field.GetOperation() {
	case filter.EQ:
		b.Builder.EQ(
			b.Field,
			marshaledField,
		)
	case filter.GT:
		b.Builder.GT(
			b.Field,
			marshaledField,
		)
	case filter.GTE:
		b.Builder.GTE(
			b.Field,
			marshaledField,
		)
	case filter.LIKE:
		str := fmt.Sprint(marshaledField)
		b.Builder.Like(
			b.Field,
			str,
			builder.I,
		)
	case filter.LT:
		b.Builder.LT(
			b.Field,
			marshaledField,
		)
	case filter.LTE:
		b.Builder.LTE(
			b.Field,
			marshaledField,
		)
	case filter.NEQ:
		b.Builder.NEQ(
			b.Field,
			marshaledField,
		)
	case filter.IN:
		b.Builder.InArray(
			b.Field,
			marshaledField,
		)
	case filter.NIN:
		b.Builder.NotInArray(
			b.Field,
			marshaledField,
		)
	}
}