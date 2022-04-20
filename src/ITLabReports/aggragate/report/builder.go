package report

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/entity/assignees"
	"github.com/RTUITLab/ITLab-Reports/entity/report"
)

type reportBuilder struct {
	report *report.Report
	assignees *assignees.Assignees
}

func NewReportBuilder() ReportBuilder {
	return &reportBuilder{}
}

type ReportBuilder interface {
	SetName(name string) ReportBuilder

	SetText(text string) ReportBuilder

	SetReporter(reporter string) ReportBuilder

	SetImplementor(implementer string) ReportBuilder

	Create() (*Report, error)
}


func (r *reportBuilder) createReportIfNil() {
	if r.report == nil {
		r.report = &report.Report{}
	}
}

func (r *reportBuilder) createAssigneesIfNil() {
	if r.assignees == nil {
		r.assignees = &assignees.Assignees{}
	}
}

func (r *reportBuilder) SetName(name string) ReportBuilder {
	r.createReportIfNil()
	r.report.Name = name

	return r
}

func (r *reportBuilder) SetText(text string) ReportBuilder {
	r.createReportIfNil()
	r.report.Text = text

	return r
}

func (r *reportBuilder) SetReporter(reporter string) ReportBuilder {
	r.createAssigneesIfNil()
	r.assignees.Reporter = reporter

	return r
}

func (r *reportBuilder) SetImplementor(implementer string) ReportBuilder {
	r.createAssigneesIfNil()
	r.assignees.Implementer = implementer

	return r
}

func (r *reportBuilder) Create() (*Report, error) {
	report := &Report{
		Report: r.report,
		Assignees: r.assignees,
	}

	report.Report.Date = time.Now().UTC().Round(time.Millisecond)

	if err := NewReportValidator().Validate(report); err != nil {
		return nil, err
	}

	return report, nil
}
