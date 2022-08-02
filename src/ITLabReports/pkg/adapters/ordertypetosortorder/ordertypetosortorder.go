package ordertypetosortorder

import (
	"github.com/0B1t322/MongoBuilder/operators/sort"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
)

func ToSortOrder(order ordertype.OrderType) sort.SortOrder {
	switch order {
	case ordertype.ASC:
		return sort.ASC()
	case ordertype.DESC:
		return sort.DESC()
	default:
		return sort.ASC()
	}
}

