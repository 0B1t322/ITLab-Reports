package aggregate

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/models/entity"
	"github.com/RTUITLab/ITLab-Reports/internal/models/valueobject"
)

type Draft struct {
	*entity.Report
	Assignees valueobject.Assignees
}

func NewDraft(
	name string,
	text string,
	reporter string,
	implementer string,
) (Draft, error) {
	draft := Draft{
		Report: &entity.Report{
			Name: name,
			Date: time.Now().UTC(),
			Text: text,
		},
		Assignees: valueobject.Assignees{
			Reporter:    reporter,
			Implementer: implementer,
		},
	}

	return draft, NewDraftValidator().Validate(draft)
}

func (d Draft) SetName(name string) {
	d.Name = name
}

func (d Draft) SetText(text string) {
	d.Text = text
}

func (d Draft) SetAssignees(reporter string, implementer string) {
	d.Assignees.Reporter = reporter
	d.Assignees.Implementer = implementer
}

func (d Draft) ToReport() Report {
	return Report{
		Report: &entity.Report{
			Name: d.Name,
			Date: d.Date,
			Text: d.Text,
		},
		Assignees: d.Assignees,
		State:     ReportStateCreated,
	}
}

func (d Draft) UserIsDraftOwner(user User) bool {
	return d.Assignees.Reporter == user.ID
}
