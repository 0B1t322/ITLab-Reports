package ordertype

type OrderType int8

const (
	ASC  = 1
	DESC = -1
)

// OrderTypeFromString get orderType in next formats:
// 	"asc", "ASC", "desc", "DESC", "1", "-1"
func OrderTypeFromString(orderType string) OrderType {
	switch orderType {
	case "asc", "ASC", "1":
		return ASC
	case "desc", "DESC", "-1":
		return DESC
	default:
		return 0
	}
}

/*
OrderTypeFromInt get orderType in next formats:
	1 - asc
	-1 - desc
*/
func OrderTypeFromInt(orderType int) OrderType {
	result := OrderType(orderType)
	switch result {
	case ASC, DESC:
		return result
	default:
		return 0
	}
}

type TypeOrderer interface {
	GetOrderType() OrderType
}