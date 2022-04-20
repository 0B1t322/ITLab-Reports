package reqresp

import (
	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
)

type GetReportReq struct {
	ID string
}

type GetReportResp struct {
	Report *report.Report
}

func (g *GetReportResp) GetReporter() string {
	return g.Report.GetReporter()
}

func (g *GetReportResp) GetImplementer() string {
	return g.Report.GetImplementer()
}