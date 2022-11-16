package sort

import (
	"strings"

	"github.com/0B1t322/RepoGen/pkg/sortorder"
	"github.com/RTUITLab/ITLab-Reports/internal/common/errors"
	"github.com/samber/lo"
)

var (
	ErrSortParamsRepeat = errors.New("You can't sort by one field twice")
)

type SortOrder int

const (
	SortOrderAsc SortOrder = iota
	SortOrderDesc
	Unknown
)

func SortOrderFromString(order string) SortOrder {
	switch order {
	case "asc", "1", "ASC":
		return SortOrderAsc
	case "desc", "-1", "DESC":
		return SortOrderDesc
	default:
		return Unknown
	}
}

// Return -1 if sortorder is uknown
func (s SortOrder) ToSortOrder() sortorder.SortOrder {
	if s == SortOrderAsc {
		return sortorder.ASC
	} else if s == SortOrderDesc {
		return sortorder.DESC
	}

	return -1
}

type SortParameter interface {
	Field() string
	Order() SortOrder
}

type SortParameters []SortParameter

// Check that sort params is not repeat
func (s SortParameters) IsFieldsRepeat() error {
	uniqs := lo.UniqBy(
		s,
		func(param SortParameter) string {
			return param.Field()
		},
	)

	if len(uniqs) != len(s) {
		return ErrSortParamsRepeat
	}

	return nil
}

type sortParameter struct {
	fieldName string
	order     SortOrder
}

func (s sortParameter) Field() string {
	return s.fieldName
}

func (s sortParameter) Order() SortOrder {
	return s.order
}

type SortMarshaler struct {
}

func (s SortMarshaler) Marshal(sort []string) (sortParams SortParameters) {
	for _, param := range sort {
		fieldAndOrder := strings.SplitN(param, ":", 2)
		if len(fieldAndOrder) != 2 {
			continue
		}

		order := SortOrderFromString(fieldAndOrder[1])
		if order == Unknown {
			continue
		}

		sortParams = append(sortParams, sortParameter{
			fieldName: fieldAndOrder[0],
			order:     order,
		})
	}

	return
}
