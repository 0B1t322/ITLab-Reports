package utils

import (
	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	"github.com/RTUITLab/ITLab-Reports/entity/assignees"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab/proto/reports/types"
	"github.com/RTUITLab/ITLab/proto/shared"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ReportToPBType(report *report.Report) *types.Report {
	return &types.Report{
		Id:        report.GetID(),
		Name:      report.GetName(),
		Text:      report.GetText(),
		Assignees: AssigneesToPBType(report.Assignees),
		Date:      timestamppb.New(report.GetDate()),
	}

}

func AssigneesToPBType(assignees *assignees.Assignees) *types.Assignees {
	return &types.Assignees{
		Implementer: assignees.Implementer,
		Reporter:    assignees.Reporter,
	}
}

func OrderTypeFromGRPC(order shared.Ordering) ordertype.OrderType {
	switch order {
	case shared.Ordering_DESCENDING:
		return ordertype.DESC
	default:
		return ordertype.ASC
	}
}
