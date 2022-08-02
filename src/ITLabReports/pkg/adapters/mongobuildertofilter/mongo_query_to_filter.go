package mongobuildertofilter

import (
	"github.com/0B1t322/MongoBuilder/operators/options"
	"github.com/0B1t322/MongoBuilder/operators/query"
	"github.com/0B1t322/MongoBuilder/utils"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"go.mongodb.org/mongo-driver/bson"
)

type BuilderQueryAdapter[T filter.FieldType] struct {
	raw bson.M

	src *bson.M

	Field string

	FieldMarshaller *FieldMarshaller[T]
}

type BuilderQueryAdapterOptions[T filter.FieldType] func (b *BuilderQueryAdapter[T])

func NewBuilderQueryAdapterOptions[T filter.FieldType]() BuilderQueryAdapterOptions[T] {
	return func (b *BuilderQueryAdapter[T]) {}
}

func (b BuilderQueryAdapterOptions[T]) WithFieldFormatter(fieldMarshaller FieldMarshaller[T]) BuilderQueryAdapterOptions[T] {
	return func(a *BuilderQueryAdapter[T]) {
		a.FieldMarshaller = &fieldMarshaller
		b(a)
	}
}

func (b BuilderQueryAdapterOptions[T]) withSetDefault() BuilderQueryAdapterOptions[T] {
	return func(a *BuilderQueryAdapter[T]) {
		if a.FieldMarshaller == nil {
			def := DefaultMarshaller[T]()
			a.FieldMarshaller = &def
		}
		b(a)
	}
}

func NewBuilderQueryAdapter[T filter.FieldType](
	b *bson.M,
	fieldName string,
	opts ...BuilderQueryAdapterOptions[T],
) filter.FilterBuilder[T] {
	builder := &BuilderQueryAdapter[T]{
		raw: bson.M{},
		src: b,
		Field: fieldName,
	}

	for _, opt := range opts {
		opt(builder)
	}

	// Set default options value
	NewBuilderQueryAdapterOptions[T]().withSetDefault()(builder)

	return builder
}

func (b *BuilderQueryAdapter[T]) MarshallField(field filter.FieldFilterer[T]) any {
	return b.FieldMarshaller.Marshaller(field)
}

func (b *BuilderQueryAdapter[T]) BuildFilterField(
	field filter.FieldFilterer[T],
) {
	marshaledField := b.MarshallField(field)
	
	switch field.GetOperation() {
	case filter.EQ:
		b.raw = query.EQField(b.Field, marshaledField)
	case filter.NEQ:
		b.raw = query.NE(b.Field, marshaledField)
	case filter.GT:
		b.raw = query.GT(b.Field, marshaledField)
	case filter.GTE:
		b.raw = query.GTE(b.Field, marshaledField)
	case filter.LT:
		b.raw = query.LT(b.Field, marshaledField)
	case filter.LTE:
		b.raw = query.LTE(b.Field, marshaledField)
	case filter.IN:
		if b.FieldMarshaller.Type == SliceID {
			b.raw = query.In(b.Field, marshaledField.(bson.A)...)
		} else {
			b.raw = query.In(b.Field, marshaledField)
		}
	case filter.NIN:
		if b.FieldMarshaller.Type == SliceID {
			b.raw = query.Nin(b.Field, marshaledField.(bson.A)...)
		} else {
			b.raw = query.Nin(b.Field, marshaledField)
		}
	case filter.LIKE:
		b.raw = query.Regex(b.Field, marshaledField.(string), options.I)
	case filter.EXIST:
		b.raw = query.Exists(b.Field, marshaledField.(bool))
	}

	*b.src = utils.MergeBsonM(*b.src, b.raw)
}	