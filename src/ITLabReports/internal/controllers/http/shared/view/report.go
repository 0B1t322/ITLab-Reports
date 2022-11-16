package view

import (
	"time"

	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
)

type ReportView struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Text      string        `json:"text"`
	Assignees AssigneesView `json:"assignees"`
	Date      time.Time     `json:"date"`
}

func ReportViewFrom(from aggregate.Report) ReportView {
	return ReportView{
		ID:        from.ID,
		Name:      from.Name,
		Text:      from.Text,
		Assignees: AssigneesViewFrom(from.Assignees),
		Date:      from.Date.UTC(),
	}
}

func ReportsViewFrom(from []aggregate.Report) []ReportView {
	var reports []ReportView

	for _, report := range from {
		reports = append(reports, ReportViewFrom(report))
	}

	return reports
}

type ReportsPageView struct {
	Reports []ReportView `json:"reports"`
}

func ReportsPageViewFrom(from []aggregate.Report) ReportsPageView {
	return ReportsPageView{
		Reports: ReportsViewFrom(from),
	}
}
