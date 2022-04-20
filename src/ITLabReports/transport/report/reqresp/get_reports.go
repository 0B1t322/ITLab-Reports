package reqresp

import (
	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
)

type GetReportsReq struct {
	Params *reportdomain.GetReportsParams
}

func (g *GetReportsReq) SetImplementerAndReporter(implementer, reporter string) {
	if g.Params == nil {
		g.Params = &reportdomain.GetReportsParams{}
	}

	if g.Params.Filter == nil {
		g.Params.Filter = &reportdomain.GetReportsFilter{}
	}

	g.Params.Filter.Or = append(
		g.Params.Filter.Or,
		&reportdomain.GetReportsFilterFieldsWithOrAnd{
			GetReportsFilterFields: reportdomain.GetReportsFilterFields{
				Reporter: &filter.FilterField[string]{
					Value: reporter,
					Operation: filter.EQ,
				},
			},
		},
		&reportdomain.GetReportsFilterFieldsWithOrAnd{
			GetReportsFilterFields: reportdomain.GetReportsFilterFields{
				Reporter: &filter.FilterField[string]{
					Value: implementer,
					Operation: filter.EQ,
				},
			},
		},
	)
}

type GetReportsResp struct {
	Reports []*report.Report
}