package views

import (
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	types "github.com/RTUITLab/ITLab/proto/reports/types"
	_ "github.com/RTUITLab/ITLab/proto/reports/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ReportFrom(report aggregate.Report) *types.Report {
	return &types.Report{
		Id:        report.ID,
		Name:      report.Name,
		Text:      report.Text,
		Assignees: AssigneesFrom(report.Assignees),
		Date:      timestamppb.New(report.Date),
	}
}
