package views

import (
	"github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"
	types "github.com/RTUITLab/ITLab/proto/reports/types"
	_ "github.com/RTUITLab/ITLab/proto/reports/v1"
)

func AssigneesFrom(assignees valueobject.Assignees) *types.Assignees {
	return &types.Assignees{
		Implementer: assignees.Implementer,
		Reporter:    assignees.Reporter,
	}
}
