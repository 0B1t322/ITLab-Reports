package reqresp

import "github.com/RTUITLab/ITLab-Reports/aggragate/report"

type CreateReportReq struct {
	Report *report.Report
}

func (c *CreateReportReq) SetImplementor(implementor string) {
	c.Report.Assignees.Implementer = implementor
}

func (c *CreateReportReq) SetReporter(reporter string) {
	c.Report.Assignees.Reporter = reporter
}

type CreateReportResp struct {
	Report *report.Report
}