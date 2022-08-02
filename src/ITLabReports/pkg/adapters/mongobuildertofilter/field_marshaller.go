package mongobuildertofilter

import (
	"fmt"
	"time"

	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FieldMarhallerType int
const (
	None FieldMarhallerType = iota
	FieldToTime
	StringID
	SliceID
)

type FieldMarshaller[T filter.FieldType] struct {
	Marshaller func(field filter.FieldFilterer[T]) any
	Type FieldMarhallerType
}

func DefaultMarshaller[T filter.FieldType]() FieldMarshaller[T]{
	return FieldMarshaller[T]{
		Marshaller: func(field filter.FieldFilterer[T]) any {return field.GetValue()},
		Type: None,
	}
}

func FieldToTimeMarshaller[T filter.FieldType]() FieldMarshaller[T] {
	return FieldMarshaller[T]{
		Marshaller: func(field filter.FieldFilterer[T]) any {
			str := fmt.Sprint(field.GetValue())
	
			time, err := time.Parse(time.RFC3339, str)
			if err != nil {
				return nil
			}
	
			return time.UTC()
		},
		Type: FieldToTime,
	}
}

func StringIdMarshaller() FieldMarshaller[string] {
	return FieldMarshaller[string]{
		Marshaller: func(field filter.FieldFilterer[string]) any {
			str := fmt.Sprint(field.GetValue())
	
			id, err := primitive.ObjectIDFromHex(str)
			if err != nil {
				return nil
			}
	
			return id
		},
		Type: FieldToTime,
	}
}

func SliceIdMarshaller() FieldMarshaller[[]string] {
	return FieldMarshaller[[]string]{
		Marshaller: func(field filter.FieldFilterer[[]string]) any {
			result := bson.A{}
			for _, elem := range field.GetValue() {
				id, err := primitive.ObjectIDFromHex(elem)
				if err != nil {
					continue
				}
				result = append(result, id)
			}
			return result
		},
		Type: SliceID,
	}
}