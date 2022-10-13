package report

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/entity/assignees"
	"github.com/RTUITLab/ITLab-Reports/entity/report"
)

type Report struct {
	Report *report.Report

	Assignees *assignees.Assignees
}

func (r Report) GetID() string {
	return r.Report.ID
}

func (r Report) GetDateString() string {
	return r.Report.Date.Format(time.RFC3339Nano)
}

func (r Report) GetDate() time.Time {
	return r.Report.Date
}

func (r Report) GetName() string {
	return r.Report.Name
}

func (r Report) GetText() string {
	return r.Report.Text
}

func (r Report) GetImplementer() string {
	return r.Assignees.Implementer
}

func (r Report) GetReporter() string {
	return r.Assignees.Reporter
}

func (r *Report) SetPaid() {
	r.Report.State = report.ReportStatePaid
}

func NewReport(
	name string,
	text string,
	reporter string,
	implementor string,
) (*Report, error) {
	report := &Report{
		Report: &report.Report{
			Name:  name,
			Text:  text,
			Date:  time.Now().UTC().Round(time.Millisecond),
			State: report.ReportStateCreated,
		},
		Assignees: &assignees.Assignees{
			Reporter:    reporter,
			Implementer: implementor,
		},
	}

	if err := NewReportValidator().Validate(report); err != nil {
		return nil, err
	}

	return report, nil
}
