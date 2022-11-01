package utils

import (
	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/entity/assignees"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func DraftToPBType(draft *report.Report) *types.Draft {
	return &types.Draft{
		Id:        draft.GetID(),
		Name:      draft.GetName(),
		Text:      draft.GetText(),
		Assignees: AssigneesToPBType(draft.Assignees),
		Date:      timestamppb.New(draft.GetDate()),
	}
}

func AssigneesToPBType(assignees *assignees.Assignees) *types.Assignees {
	return &types.Assignees{
		Implementer: assignees.Implementer,
		Reporter:    assignees.Reporter,
	}
}
