package reqresp

import (
	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/aggragate/report"
)

type UpdateReportReq struct {
	ID string
	Params reportdomain.UpdateReportParams
}

type UpdateReportResp struct {
	Report *report.Report
}