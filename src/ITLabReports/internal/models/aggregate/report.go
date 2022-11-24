package aggregate

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/models/entity"
	"github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"
)

type ReportState string

const (
	ReportStateCreated ReportState = "created"
	ReportStatePaid    ReportState = "paid"
)

type Report struct {
	*entity.Report
	Assignees valueobject.Assignees
	State     ReportState
}

func NewReport(
	name string,
	text string,
	reporter string,
	implementer string,
) (Report, error) {
	report := Report{
		Report: &entity.Report{
			Name: name,
			Date: time.Now().UTC(),
			Text: text,
		},
		Assignees: valueobject.Assignees{
			Reporter:    reporter,
			Implementer: implementer,
		},
		State: ReportStateCreated,
	}

	return report, NewReportValidator().Validate(report)
}

func (r Report) SetName(name string) {
	r.Name = name
}

func (r Report) SetText(text string) {
	r.Text = text
}

func (r Report) SetAssignees(reporter string, implementer string) {
	r.Assignees.Reporter = reporter
	r.Assignees.Implementer = implementer
}

func (r *Report) SetPaid() {
	r.State = ReportStatePaid
}

func (r Report) UserIsReportOwner(user User) bool {
	return r.Assignees.Reporter == user.ID
}

func (r Report) ReportAboutUser(user User) bool {
	return r.Assignees.Implementer == user.ID
}
