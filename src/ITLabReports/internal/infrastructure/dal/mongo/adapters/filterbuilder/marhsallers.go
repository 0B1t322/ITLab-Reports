package filterbuilder

import (
	"fmt"
	"time"

	"github.com/0B1t322/RepoGen/pkg/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		result := bson.A{}
		for _, elem := range field.GetValue() {
			id, err := primitive.ObjectIDFromHex(elem)
			if err != nil {
				continue
			}
			result = append(result, id)
		}
		return result
	}
}

func StringSliceMarshaller() FieldMarshaller[[]string] {
	return func(field filter.FieldFilterer[[]string]) any {
		result := bson.A{}
		for _, elem := range field.GetValue() {
			result = append(result, elem)
		}

		return result
	}
}
