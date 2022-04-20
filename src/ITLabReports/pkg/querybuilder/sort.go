package querybuilder

import "github.com/RTUITLab/ITLab-Reports/pkg/ordertype"

type SortQuery struct {
	Field     string
	OrderType ordertype.TypeOrderer
}
