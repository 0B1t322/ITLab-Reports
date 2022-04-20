package ordertypetomongo

import "github.com/RTUITLab/ITLab-Reports/pkg/ordertype"

func ToMongoOrderType(order ordertype.OrderType) int {
	switch order {
	case ordertype.ASC:
		return 1
	case ordertype.DESC:
		return -1
	default:
		return 0
	}
}