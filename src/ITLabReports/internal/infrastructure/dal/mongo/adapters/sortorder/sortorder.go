package sortorder

import (
	"github.com/0B1t322/MongoBuilder/operators/sort"
	"github.com/0B1t322/RepoGen/pkg/sortorder"
)

type SortOrder = sortorder.SortOrder

func ToSortOrder(order SortOrder) sort.SortOrder {
	switch order {
	case sortorder.ASC:
		return sort.ASC()
	case sortorder.DESC:
		return sort.DESC()
	default:
		return sort.ASC()
	}
}
