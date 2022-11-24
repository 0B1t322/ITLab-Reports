package utils

import (
	"github.com/0B1t322/RepoGen/pkg/sortorder"
	"github.com/RTUITLab/ITLab/proto/shared"
)

func ProtoSortOrderTo(order shared.Ordering) sortorder.SortOrder {
	switch order {
	case shared.Ordering_ASCENDING:
		return sortorder.ASC
	default:
		return sortorder.DESC
	}
}
