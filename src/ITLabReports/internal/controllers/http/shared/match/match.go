package match

import (
	"strings"

	"github.com/samber/lo"
)

type IMatchParam[FieldType ~string, Value any] interface {
	GetField() FieldType
	GetValue() Value
}

type IMatchDesc[ParamNameType ~string, ValueType any] interface {
	Param() ParamNameType
	ValueFromString(string) (ValueType, error)
}

type matchDesc[FieldType ~string, Value any] struct {
	ParamName           FieldType
	ValueFromStringFunc func(string) (Value, error)
}

func (m matchDesc[FieldType, Value]) Param() FieldType {
	return m.ParamName
}

func (m matchDesc[FieldType, Value]) ValueFromString(s string) (Value, error) {
	return m.ValueFromStringFunc(s)
}

func NewMatchDesc[ParamNameType ~string, ValueType any](
	ParamFieldName ParamNameType,
	ValueFromString func(string) (ValueType, error),
) IMatchDesc[ParamNameType, ValueType] {
	return &matchDesc[ParamNameType, ValueType]{
		ParamName:           ParamFieldName,
		ValueFromStringFunc: ValueFromString,
	}
}

type MatchParameter interface {
	Field() string
	Value() string
}

type matchParameter struct {
	fieldName string
	value     string
}

func (m matchParameter) Field() string {
	return m.fieldName
}

func (m matchParameter) Value() string {
	return m.value
}

type MatchParameters []MatchParameter

type MatchMarshaler struct {
}

func (m MatchMarshaler) Marshal(match []string) (matchParams MatchParameters) {
	lo.ForEach(
		match,
		func(param string, _ int) {
			fieldAndValue := strings.SplitN(param, ":", 2)
			if len(fieldAndValue) != 2 {
				return
			}

			matchParams = append(matchParams, matchParameter{
				fieldName: fieldAndValue[0],
				value:     fieldAndValue[1],
			})
		},
	)
	return matchParams
}
